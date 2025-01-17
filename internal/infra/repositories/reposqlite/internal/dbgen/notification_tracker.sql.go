// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: notification_tracker.sql

package dbgen

import (
	"context"
)

const GetNotificationByEvent = `-- name: GetNotificationByEvent :many
SELECT id, event_id, event_sha256, status, notified_at
FROM notified_events
WHERE notified_events.event_id = ?
`

type GetNotificationByEventRow struct {
	ID          int64       `db:"id"`
	EventID     int64       `db:"event_id"`
	EventSha256 string      `db:"event_sha256"`
	Status      string      `db:"status"`
	NotifiedAt  interface{} `db:"notified_at"`
}

func (q *Queries) GetNotificationByEvent(ctx context.Context, eventID int64) ([]GetNotificationByEventRow, error) {
	rows, err := q.db.QueryContext(ctx, GetNotificationByEvent, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetNotificationByEventRow{}
	for rows.Next() {
		var i GetNotificationByEventRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.EventSha256,
			&i.Status,
			&i.NotifiedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const MarkEventAsNotified = `-- name: MarkEventAsNotified :one
INSERT INTO notified_events (event_id, event_sha256, status, notified_at)
VALUES (?, ?, ?, ?)
ON CONFLICT (event_sha256) DO UPDATE SET event_sha256=excluded.event_sha256,
                                         status=excluded.status,
                                         notified_at=excluded.notified_at
RETURNING id
`

type MarkEventAsNotifiedParams struct {
	EventID     int64       `db:"event_id"`
	EventSha256 string      `db:"event_sha256"`
	Status      string      `db:"status"`
	NotifiedAt  interface{} `db:"notified_at"`
}

func (q *Queries) MarkEventAsNotified(ctx context.Context, arg MarkEventAsNotifiedParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, MarkEventAsNotified,
		arg.EventID,
		arg.EventSha256,
		arg.Status,
		arg.NotifiedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}
