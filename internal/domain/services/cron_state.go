package services

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

type (
	eventObserver interface {
		OnEvent(event entities.OlympicEvent)
	}
	jobContext struct {
		ID      uuid.UUID
		startAt time.Time
		endAt   time.Time
	}

	cronState struct {
		cronScheduler       gocron.Scheduler
		jobIDs              map[string]jobContext
		registeredObservers []eventObserver
		mutex               sync.Mutex
	}
)

func (cs *cronState) retrieveJobID(key string) (jobContext, bool) {
	if cs.jobIDs == nil {
		return jobContext{ID: uuid.Nil}, false
	}

	c, ok := cs.jobIDs[key]
	return c, ok
}

func (cs *cronState) registerJobInfo(key string, jCtx jobContext) {
	cs.jobIDs[key] = jCtx
}

func (cs *cronState) RegisterObserver(observer eventObserver) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.registeredObservers = append(cs.registeredObservers, observer)
}

func (cs *cronState) taskExecution(event entities.OlympicEvent) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for _, observer := range cs.registeredObservers {
		observer.OnEvent(event)
	}
}
