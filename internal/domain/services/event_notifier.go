package services

import (
	"errors"
	"log/slog"
	"sync"
	"time"

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
		usecases.CanNotifyRepository
	}
)

type (
	CancelChannel    chan struct{}
	notifierUseCases struct {
		usecases.CanNotifyUseCase
		usecases.FetcherCacheUseCase
	}
	EventNotifier struct {
		cancelChan      CancelChannel
		checkInterval   time.Duration
		olympicsEndDate time.Time
		useCases        notifierUseCases
		repo            EventNotifierRepository
		mutex           sync.Mutex
		cronState
	}
)

func NewEventNotifier(
	cancelChan CancelChannel, cacheDuration time.Duration,
	repo EventNotifierRepository,
	fetcherUseCase usecases.FetcherCacheUseCase, canNotifyUseCase usecases.CanNotifyUseCase,
) (en *EventNotifier, err error) {
	en = &EventNotifier{
		cancelChan:    cancelChan,
		checkInterval: cacheDuration,
		repo:          repo,
		useCases: notifierUseCases{
			CanNotifyUseCase:    canNotifyUseCase,
			FetcherCacheUseCase: fetcherUseCase,
		},
		olympicsEndDate: time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC),
	}

	return
}

func (en *EventNotifier) MainLoop() error {
	go en.fetcherThread()

	ticker := time.NewTicker(en.checkInterval >> 1)
	defer ticker.Stop()

	// Do a first check before running the timer
	if err := en.updateDisciplines(); err != nil {
		return err
	}
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

func (en *EventNotifier) updateDisciplines() error {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	if err := en.useCases.FetchDisciplines(); err != nil {
		return err
	}
	return nil
}

func (en *EventNotifier) fetchRemainingDays(from time.Time) {
	en.mutex.Lock() // Prevent sqlite multi access
	defer en.mutex.Unlock()

	startDate := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.Add(48 * time.Hour)
	if endDate.After(en.olympicsEndDate) {
		endDate = en.olympicsEndDate
	}

	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		slog.Info("Start to fetch and save data for event", slog.Time("date", date))
		if err := en.useCases.FetcherCacheUseCase.FetchDay(date); err != nil {
			slog.Error("Error fetching data from day", slog.String("error", err.Error()))
		}
	}
}

func (en *EventNotifier) fetcherThread() {
	// Run one time since beginning
	en.fetchRemainingDays(time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC))
	ticker := time.NewTicker(en.checkInterval)
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

func (en *EventNotifier) taskExecution(event entities.OlympicEvent) {
	// Use cron state to trigger observers
	notifyStatus := entities.NotificationStatusSent
	if !en.cronState.taskExecution(event) {
		notifyStatus = entities.NotificationStatusFailed
	}

	// Update status on database
	newNotification := entities.Notification{
		EventID:       event.ID,
		Status:        notifyStatus,
		EventChecksum: event.SHAIdentifier(),
		NotifiedAt:    time.Now(),
	}
	if err := en.repo.RegisterNotification(newNotification); err != nil {
		slog.Error(
			"Error registering notification",
			slog.String("error", err.Error()),
		)
	}
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

		eventKey, checkErr := en.useCases.ShouldNotify(event)
		if checkErr != nil || eventKey == "" {
			errList = append(errList, checkErr)
			continue
		}

		// Check again for the 20min notification
		/*startDiff := utils.AbsoluteNum(event.StartAt.Sub(time.Now()))
		if startDiff <= 20*time.Minute {
			en.taskExecution(event)
		}*/
		en.taskExecution(event)

		slog.Info(
			"Job for notification sent",
			slog.String("eventHash", event.SHAIdentifier()),
			slog.Time("startTime", event.StartAt),
		)

	}

	return errors.Join(errList...)
}
