package reposqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) CheckSentNotifications(
	eventID entities.Identifier,
	eventChecksum string,
) (entities.Notification, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbQuery := r.queries
	foundNotifications, err := dbQuery.GetNotificationByEvent(ctx, int64(eventID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entities.Notification{}, err
	}

	for _, notification := range foundNotifications {
		if notification.EventSha256 == eventChecksum {
			foundNotification := entities.Notification{
				ID:            entities.Identifier(notification.ID),
				EventID:       eventID,
				Status:        entities.NotificationStatus(notification.Status),
				EventChecksum: notification.EventSha256,
			}
			foundNotification.NotifiedAt, _ = notification.NotifiedAt.(time.Time)
			return foundNotification, nil
		}
	}

	return entities.Notification{}, nil
}

func (r RepoSQLite) RegisterNotification(notification entities.Notification) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbQuery := r.queries

	var notifiedAt *time.Time
	if !notification.NotifiedAt.IsZero() {
		notifiedAt = &notification.NotifiedAt
	}
	_, err := dbQuery.MarkEventAsNotified(
		ctx, dbgen.MarkEventAsNotifiedParams{
			EventID:     int64(notification.EventID),
			EventSha256: notification.EventChecksum,
			Status:      string(notification.Status),
			NotifiedAt:  notifiedAt,
		},
	)

	return err
}
