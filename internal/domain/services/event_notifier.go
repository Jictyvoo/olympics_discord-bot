package services

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

type (
	EventLoader interface {
		LoadDayEvents(now time.Time) ([]entities.OlympicEvent, error)
		// LoadCompetitorsFromEvent(event entities.OlympicEvent) ([]entities.OlympicCompetitors, error)
	}
	EventNotifier struct {
		cancelChan     chan struct{}
		cacheDuration  time.Duration
		fetcherUseCase usecases.FetcherCacheUseCase
		repo           EventLoader
		mutex          sync.Mutex
		cronState
	}
)

func NewEventNotifier(
	cancelChan chan struct{}, cacheDuration time.Duration,
	fetcherUseCase usecases.FetcherCacheUseCase,
) (en *EventNotifier, err error) {
	en = &EventNotifier{
		cancelChan:     cancelChan,
		cacheDuration:  cacheDuration,
		fetcherUseCase: fetcherUseCase,
	}
	en.cronScheduler, err = gocron.NewScheduler()
	en.jobIDs = make(map[string]jobContext, 100)

	return
}

func (en *EventNotifier) fetchRemainingDays() {
	en.mutex.Lock() // Prevent sqlite multi access
	defer en.mutex.Unlock()

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC)
	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		slog.Debug("Start to fetch again for event data", slog.Time("date", date))
		if err := en.fetcherUseCase.Run(date); err != nil {
			slog.Error("Error fetching data from day", slog.String("error", err.Error()))
		}
	}
}

func (en *EventNotifier) fetcherThread() {
	ticker := time.NewTicker(en.cacheDuration)
	for {
		select {
		case _, _ = <-en.cancelChan:
			ticker.Stop()
			return
		case <-ticker.C:
			en.fetchRemainingDays()
		}
	}
}

func (en *EventNotifier) Start() error {
	go en.fetcherThread()

	ticker := time.NewTicker(en.cacheDuration << 1)
	defer ticker.Stop()

	for range ticker.C {
		if err := en.checkUpdateJobs(); err != nil {
			return err
		}
	}

	return nil
}

func (en *EventNotifier) manageEventJob(event entities.OlympicEvent) (err error) {
	eventKey := event.SHAIdentifier()
	jobCtx, found := en.cronState.retrieveJobID(eventKey)
	shouldInsertNew := true
	if found && jobCtx.ID != uuid.Nil {
		// Check if it has some update
		shouldInsertNew = event.StartAt.Compare(jobCtx.startAt) == 0 &&
			event.EndAt.Compare(jobCtx.endAt) == 0

		if !shouldInsertNew {
			err = en.cronScheduler.RemoveJob(jobCtx.ID)
			if err != nil {
				slog.Error("Error removing job", slog.String("error", err.Error()))
				return err
			}
		}
	}

	if shouldInsertNew {
		newJob, insertErr := en.cronScheduler.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(event.StartAt),
			),
			gocron.NewTask(en.cronState.taskExecution, event),
			gocron.WithStopAt(gocron.WithStopDateTime(event.EndAt)),
			gocron.WithLimitedRuns(1),
			gocron.WithSingletonMode(gocron.LimitModeWait),
		)

		if insertErr != nil {
			slog.Error("Error creating job", slog.String("error", insertErr.Error()))
			return insertErr
		}
		en.cronState.registerJobInfo(
			eventKey, jobContext{
				ID:      newJob.ID(),
				startAt: event.StartAt,
				endAt:   event.EndAt,
			},
		)
	}

	return
}

func (en *EventNotifier) checkUpdateJobs() error {
	en.mutex.Lock() // Prevent sqlite multi access
	defer en.mutex.Unlock()

	// Fetch all events from the day
	dayEvents, err := en.repo.LoadDayEvents(time.Now())
	if err != nil {
		slog.Error("Error loading day events", slog.String("error", err.Error()))
		return err
	}

	errList := make([]error, 0, len(dayEvents))
	for _, event := range dayEvents {
		// event.Competitors, err = en.repo.LoadCompetitorsFromEvent(event)

		if err = en.manageEventJob(event); err != nil {
			errList = append(errList, err)
		}
	}

	return errors.Join(errList...)
}
