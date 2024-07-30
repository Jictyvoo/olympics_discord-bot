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
		LoadEvent(id entities.Identifier) (entities.OlympicEvent, error)
		// LoadCompetitorsFromEvent(event entities.OlympicEvent) ([]entities.OlympicCompetitors, error)
	}

	EventNotifierRepository interface {
		EventLoader
		CheckSentNotifications(
			eventID entities.Identifier, eventChecksum string,
		) (entities.Notification, error)
		RegisterNotification(notification entities.Notification) error
	}
)

type (
	CancelChannel chan struct{}
	EventNotifier struct {
		cancelChan     CancelChannel
		cacheDuration  time.Duration
		fetcherUseCase usecases.FetcherCacheUseCase
		repo           EventNotifierRepository
		mutex          sync.Mutex
		cronState
	}
)

func NewEventNotifier(
	cancelChan CancelChannel, cacheDuration time.Duration,
	repo EventNotifierRepository, fetcherUseCase usecases.FetcherCacheUseCase,
) (en *EventNotifier, err error) {
	en = &EventNotifier{
		cancelChan:     cancelChan,
		cacheDuration:  cacheDuration,
		repo:           repo,
		fetcherUseCase: fetcherUseCase,
	}
	en.cronScheduler, err = gocron.NewScheduler()
	en.jobIDs = make(map[string]jobContext, 100)

	return
}

func (en *EventNotifier) fetchRemainingDays(from time.Time) {
	en.mutex.Lock() // Prevent sqlite multi access
	defer en.mutex.Unlock()

	startDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC)
	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		slog.Info("Start to fetch and save data for event", slog.Time("date", date))
		if err := en.fetcherUseCase.Run(date); err != nil {
			slog.Error("Error fetching data from day", slog.String("error", err.Error()))
		}
	}
}

func (en *EventNotifier) fetcherThread() {
	// Run one time since beginning
	en.fetchRemainingDays(time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC))
	ticker := time.NewTicker(en.cacheDuration)
	for {
		select {
		case _, _ = <-en.cancelChan:
			ticker.Stop()
			return
		case <-ticker.C:
			en.fetchRemainingDays(time.Now())
		}
	}
}

func (en *EventNotifier) Start() {
	en.cronScheduler.Start()
}

func (en *EventNotifier) MainLoop() error {
	go en.fetcherThread()

	ticker := time.NewTicker(en.cacheDuration >> 1)
	defer ticker.Stop()

	// Do a first check before running the timer
	if err := en.checkUpdateJobs(); err != nil {
		return err
	}

	for {
		select {
		case _, _ = <-en.cancelChan:
			return nil
		case _, _ = <-ticker.C:
			if err := en.checkUpdateJobs(); err != nil {
				return err
			}
		}
	}
}

func (en *EventNotifier) taskExecution(event entities.OlympicEvent) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	updatedEvent, err := en.repo.LoadEvent(event.ID)
	if err != nil {
		slog.Error(
			"Error loading event",
			slog.String("error", err.Error()),
		)
		updatedEvent = event
	}
	en.cronState.taskExecution(updatedEvent)
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

	startTime := event.StartAt.Add(-20 * time.Minute)
	now := time.Now()
	if shouldInsertNew && event.EndAt.After(now) {
		if startTime.Before(now) {
			startTime = now.Add(10 * time.Second)
		}
		newJob, insertErr := en.cronScheduler.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(startTime),
			),
			gocron.NewTask(en.taskExecution, event),
			// gocron.WithStopAt(gocron.WithStopDateTime(event.EndAt)),
			gocron.WithLimitedRuns(1),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
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
