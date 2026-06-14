package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type NotificationRepo struct{ *repoMySQL }

func NewNotificationRepo(base *repoMySQL) NotificationRepo { return NotificationRepo{base} }

func (r NotificationRepo) UpsertNotification(n eventcore.Notification) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertNotification(qctx, dbgen.UpsertNotificationParams{
		ID:        n.ID.Bytes(),
		AlertID:   n.AlertID.Bytes(),
		ChannelID: mapper.NSStr(n.ChannelID),
		MessageID: mapper.NSStr(n.MessageID),
		Status:    string(n.Status),
		Checksum:  mapper.NSStr(n.Checksum),
		SentAt:    mapper.NSTime(n.SentAt),
	})
}

func (r NotificationRepo) GetNotificationByChecksum(
	checksum string,
) (eventcore.Notification, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetNotificationByChecksum(qctx, mapper.NSStr(checksum))
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

func rowToNotification(row dbgen.Notification) eventcore.Notification {
	return eventcore.Notification{
		ID:        mapper.IDFromBytes(row.ID),
		AlertID:   mapper.IDFromBytes(row.AlertID),
		ChannelID: row.ChannelID.String,
		MessageID: row.MessageID.String,
		Status:    eventcore.ParseNotificationStatus(row.Status),
		Checksum:  row.Checksum.String,
		SentAt:    row.SentAt.Time.UTC(),
	}
}
