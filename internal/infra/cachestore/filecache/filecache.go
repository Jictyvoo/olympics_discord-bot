package filecache

import (
	"context"
	"io"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	cacheDirPerm           os.FileMode = 0o750
	cacheFilePerm          os.FileMode = 0o600
	defaultCacheTTLMinutes             = 4
)

func createDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, cacheDirPerm)
	}
	return nil
}

type entry struct {
	key        string
	identifier string
	content    []byte
}

type FileCache struct {
	rootPath   string
	defaultTTL time.Duration
	loaded     map[string]entry
}

func New(rootPath string, defaultTTL time.Duration) (*FileCache, error) {
	if err := createDirIfNotExist(rootPath); err != nil {
		return nil, err
	}
	if defaultTTL == 0 {
		defaultTTL = defaultCacheTTLMinutes * time.Minute
	}
	c := &FileCache{rootPath: rootPath, defaultTTL: defaultTTL}
	var err error
	c.loaded, err = loadDir(rootPath)
	return c, err
}

func loadDir(dir string) (map[string]entry, error) {
	result := make(map[string]entry)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		path := filepath.Join(dir, f.Name())
		if f.IsDir() {
			sub, subErr := loadDir(path)
			if subErr == nil {
				maps.Copy(result, sub)
			}
			continue
		}
		data, readErr := os.ReadFile(filepath.Clean(path))
		if readErr != nil {
			return nil, readErr
		}
		result[f.Name()] = entry{key: f.Name(), identifier: path, content: data}
	}
	return result, nil
}

func (c *FileCache) Read(_ context.Context, key string) ([]byte, bool, error) {
	filename := filepath.Join(c.rootPath, c.subDir(0), key)
	if e, ok := c.loaded[key]; ok && e.identifier == filename {
		return e.content, true, nil
	}
	data, err := os.ReadFile(filepath.Clean(filename))
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	c.loaded[key] = entry{key: key, identifier: filename, content: data}
	return data, true, nil
}

func (c *FileCache) Write(_ context.Context, key string, value []byte, ttl time.Duration) error {
	subDir := filepath.Join(c.rootPath, c.subDir(ttl))
	filename := filepath.Join(subDir, key)
	c.loaded[key] = entry{key: key, identifier: filename, content: value}
	if err := createDirIfNotExist(subDir); err != nil {
		return err
	}
	f, err := os.OpenFile(
		filepath.Clean(filename),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		cacheFilePerm,
	)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = io.WriteString(f, "")
	if err != nil {
		return err
	}
	_, err = f.Write(value)
	return err
}

func (c *FileCache) Delete(_ context.Context, key string) error {
	e, ok := c.loaded[key]
	if !ok {
		return nil
	}
	delete(c.loaded, key)
	return os.Remove(e.identifier)
}

func (c *FileCache) subDir(ttl time.Duration) string {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	now := time.Now()
	base := now.Format("20060102")
	day := now.Day()
	hour := now.Hour()
	var div int
	switch {
	case ttl < time.Hour:
		div = now.Minute() / int(ttl.Minutes())
	case ttl < 24*time.Hour:
		div = now.Hour() / int(ttl.Hours())
		hour = 0
	}
	return base + ttl.String() + "#" + strconv.Itoa(
		day,
	) + "@" + strconv.Itoa(
		hour,
	) + "__" + strconv.Itoa(
		div,
	)
}
