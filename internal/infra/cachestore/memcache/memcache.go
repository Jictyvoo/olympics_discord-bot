package memcache

import (
	"context"
	"sync"
	"time"
)

const (
	defaultTTLMinutes = 5
	janitorTTLDivisor = 2
)

type entry struct {
	value []byte
	exp   time.Time
}

type MemCache struct {
	mu          sync.RWMutex
	data        map[string]entry
	defaultTTL  time.Duration
	stopJanitor chan struct{}
}

func New(defaultTTL time.Duration) *MemCache {
	if defaultTTL == 0 {
		defaultTTL = defaultTTLMinutes * time.Minute
	}
	c := &MemCache{
		data:        make(map[string]entry),
		defaultTTL:  defaultTTL,
		stopJanitor: make(chan struct{}),
	}
	go c.janitor(defaultTTL / janitorTTLDivisor)
	return c
}

func (c *MemCache) Close() {
	close(c.stopJanitor)
}

func (c *MemCache) Read(_ context.Context, key string) ([]byte, bool, error) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.exp) {
		return nil, false, nil
	}
	return e.value, true, nil
}

func (c *MemCache) Write(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	c.mu.Lock()
	c.data[key] = entry{value: value, exp: time.Now().Add(ttl)}
	c.mu.Unlock()
	return nil
}

func (c *MemCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return nil
}

func (c *MemCache) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.evict()
		case <-c.stopJanitor:
			return
		}
	}
}

func (c *MemCache) evict() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.data {
		if now.After(e.exp) {
			delete(c.data, k)
		}
	}
}
