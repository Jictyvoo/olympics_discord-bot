package services

import (
	"sync"
	"sync/atomic"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

type (
	EventObserver interface {
		OnEvent(event entities.OlympicEvent)
	}

	cronState struct {
		registeredObservers []EventObserver
		mutex               sync.Mutex
		started             atomic.Bool
	}
)

func (cs *cronState) Start() {
	cs.started.Store(true)
}

func (cs *cronState) RegisterObserver(observer EventObserver) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.registeredObservers = append(cs.registeredObservers, observer)
}

func (cs *cronState) taskExecution(event entities.OlympicEvent) bool {
	if !cs.started.Load() {
		return false
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for _, observer := range cs.registeredObservers {
		observer.OnEvent(event)
	}

	return true
}
