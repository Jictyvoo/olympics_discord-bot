package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type NotificationRepo struct{ *repoSQLite }

func NewNotificationRepo(base *repoSQLite) NotificationRepo { return NotificationRepo{base} }

func (r NotificationRepo) UpsertNotification(n eventcore.Notification) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	var sentAt any
	if !n.SentAt.IsZero() {
		sentAt = n.SentAt.UTC()
	}
	return r.Queries().UpsertNotification(qctx, dbgen.UpsertNotificationParams{
		ID:        n.ID.Bytes(),
		AlertID:   n.AlertID.Bytes(),
		ChannelID: mapper.OptString(n.ChannelID),
		MessageID: mapper.OptString(n.MessageID),
		Status:    string(n.Status),
		Checksum:  mapper.OptString(n.Checksum),
		SentAt:    sentAt,
	})
}

func (r NotificationRepo) GetNotificationByChecksum(
	checksum string,
) (eventcore.Notification, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetNotificationByChecksum(qctx, mapper.OptString(checksum))
	if err != nil {
		return eventcore.Notification{}, err
	}
	return rowToNotification(row), nil
}

func (r NotificationRepo) UpdateNotificationStatus(
	id eventcore.CanonicalID,
	status eventcore.NotificationStatus,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpdateNotificationStatus(qctx, dbgen.UpdateNotificationStatusParams{
		Status: string(status),
		ID:     id.Bytes(),
	})
}

// GetLatestSentNotificationByAlert returns the most recent sent notification
// (carrying a Discord message id) for a fixture, so it can be edited in place.
func (r NotificationRepo) GetLatestSentNotificationByAlert(
	alertID eventcore.CanonicalID,
) (eventcore.Notification, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetLatestSentNotificationByAlert(qctx, alertID.Bytes())
	if err != nil {
		return eventcore.Notification{}, err
	}
	return rowToNotification(row), nil
}

func rowToNotification(row dbgen.Notification) eventcore.Notification {
	return eventcore.Notification{
		ID:        mapper.IDFromBytes(row.ID),
		AlertID:   mapper.IDFromBytes(row.AlertID),
		ChannelID: mapper.NullStr(row.ChannelID),
		MessageID: mapper.NullStr(row.MessageID),
		Status:    eventcore.ParseNotificationStatus(row.Status),
		Checksum:  mapper.NullStr(row.Checksum),
		SentAt:    mapper.TimeOrZero(row.SentAt),
	}
}
