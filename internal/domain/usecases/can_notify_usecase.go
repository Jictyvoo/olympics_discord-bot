package usecases

import (
	"log/slog"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

//go:generate mockgen -destination=../mocks/canotify_repository_mock.go -package=mocks github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases CanNotifyRepository

type CanNotifyRepository interface {
	CheckSentNotifications(
		eventID entities.Identifier, eventChecksum string,
	) (entities.Notification, error)
	RegisterNotification(notification entities.Notification) error
}

type CanNotifyUseCase struct {
	repo            CanNotifyRepository
	allowedTimeDiff time.Duration
	timeNow         func() time.Time
}

func NewCanNotifyUseCase(
	repo CanNotifyRepository,
) CanNotifyUseCase {
	allowedTimeDiff := (3 * time.Hour) + (30 * time.Minute)
	return CanNotifyUseCase{repo: repo, allowedTimeDiff: allowedTimeDiff, timeNow: time.Now}
}

func (uc CanNotifyUseCase) timeDiffAllowed(event entities.OlympicEvent) bool {
	now := uc.timeNow()
	startDiff := utils.AbsoluteNum(event.StartAt.Sub(now))
	endDiff := utils.AbsoluteNum(event.EndAt.Sub(now))

	// Check if it is inside the correct element
	return endDiff+startDiff <= (uc.allowedTimeDiff << 1)
}

func (uc CanNotifyUseCase) ShouldNotify(
	event entities.OlympicEvent,
) (validatedKey string, err error) {
	// Remove ongoing results from event to prevent sending multiple ongoing notifications
	if event.Status != entities.StatusFinished &&
		uc.timeNow().Before(event.EndAt.Add(uc.allowedTimeDiff>>1)) {
		event.ResultPerCompetitor = map[string]entities.Results{}
	}

	(*entities.OlympicEvent).Normalize(&event)
	eventKey := event.SHAIdentifier()
	validatedKey = eventKey
	// Check if it exists on database
	notificationRegister, _ := uc.repo.CheckSentNotifications(event.ID, eventKey)
	// Check if it has the pending status
	if notificationRegister.Status != "" &&
		notificationRegister.Status != entities.NotificationStatusPending &&
		notificationRegister.Status != entities.NotificationStatusFailed {
		slog.Warn(
			"Ignoring event, because it already has a notification registered",
			slog.Any("notification", notificationRegister),
		)
		return "", nil
	}

	// Liberate for next checks
	notificationStatus := entities.NotificationStatusPending
	if !uc.timeDiffAllowed(event) {
		validatedKey = ""

		// Only mark to skip or cancel if it was already done
		if event.EndAt.Before(uc.timeNow().Add(-uc.allowedTimeDiff >> 1)) {
			notificationStatus = entities.NotificationStatusSkipped

			// Check if it exists on database
			if notificationRegister.Status != "" {
				notificationStatus = entities.NotificationStatusCancelled
			}
		}

		slog.Warn(
			"Event can't be send",
			slog.String("sessionCode", event.SessionCode),
			slog.String("eventKey", eventKey),
			slog.Time("startAt", event.StartAt),
			slog.Time("endAt", event.EndAt),
			slog.String("status", string(notificationStatus)),
		)
	}

	err = uc.repo.RegisterNotification(
		entities.Notification{
			EventID:       event.ID,
			Status:        notificationStatus,
			EventChecksum: eventKey,
		},
	)

	return
}
