-- name: UpsertNotification :exec
INSERT INTO notifications (id, alert_id, channel_id, message_id, status, checksum, sent_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    status     = excluded.status,
    message_id = excluded.message_id,
    sent_at    = excluded.sent_at,
    updated_at = DATETIME('now');

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = ? LIMIT 1;

-- name: GetNotificationByChecksum :one
SELECT * FROM notifications WHERE checksum = ? LIMIT 1;

-- name: UpdateNotificationStatus :exec
UPDATE notifications SET status = ?, updated_at = DATETIME('now') WHERE id = ?;

-- name: GetLatestSentNotificationByAlert :one
SELECT * FROM notifications
WHERE alert_id = ? AND status = 'sent' AND message_id IS NOT NULL AND message_id <> ''
ORDER BY created_at DESC LIMIT 1;
