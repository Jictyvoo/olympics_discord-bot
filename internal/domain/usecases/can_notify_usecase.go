package usecases

import (
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

type CanNotifyRepository interface {
	CheckSentNotifications(
		eventID entities.Identifier, eventChecksum string,
	) (entities.Notification, error)
	RegisterNotification(notification entities.Notification) error
}

type CanNotifyUseCase struct {
	repo            CanNotifyRepository
	allowedTimeDiff time.Duration
}

func NewCanNotifyUseCase(
	repo CanNotifyRepository,
) CanNotifyUseCase {
	allowedTimeDiff := (2 * time.Hour) + (30 * time.Minute)
	return CanNotifyUseCase{repo: repo, allowedTimeDiff: allowedTimeDiff}
}

func (uc CanNotifyUseCase) timeDiffAllowed(event entities.OlympicEvent) bool {
	now := time.Now()
	startDiff := utils.AbsoluteNum(event.StartAt.Sub(now))
	endDiff := utils.AbsoluteNum(event.EndAt.Sub(now))

	// Check if it is inside the correct element
	return endDiff+startDiff < (uc.allowedTimeDiff << 1)
}

func (uc CanNotifyUseCase) ShouldNotify(event entities.OlympicEvent) (string, error) {
	eventKey := event.SHAIdentifier()
	// Check if it exists on database
	notificationRegister, err := uc.repo.CheckSentNotifications(event.ID, eventKey)
	if err == nil && notificationRegister.ID != 0 {
		// Check if it has the pending status
		if notificationRegister.Status != entities.NotificationStatusPending &&
			notificationRegister.Status != entities.NotificationStatusFailed {
			return "", nil
		}
	}

	// Liberate for next checks
	notificationStatus := entities.NotificationStatusPending
	if !uc.timeDiffAllowed(event) {
		eventKey = ""
		// Check if it exists on database
		//goland:noinspection GoDfaErrorMayBeNotNil
		if notificationRegister.Status != "" {
			notificationStatus = entities.NotificationStatusSkipped
		}
		notificationStatus = entities.NotificationStatusCancelled
	}

	err = uc.repo.RegisterNotification(
		entities.Notification{
			EventID:       event.ID,
			Status:        notificationStatus,
			EventChecksum: event.SHAIdentifier(),
		},
	)

	return eventKey, err
}
