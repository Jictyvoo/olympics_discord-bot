-- name: GetNotificationByEvent :many
SELECT id, event_id, event_sha256, status, notified_at
FROM notified_events
WHERE notified_events.event_id = ?;


-- name: MarkEventAsNotified :one
INSERT INTO notified_events (event_id, event_sha256, status, notified_at)
VALUES (?, ?, ?, ?)
ON CONFLICT (event_sha256) DO UPDATE SET event_sha256=excluded.event_sha256,
                                         status=excluded.status,
                                         notified_at=excluded.notified_at
RETURNING id;
