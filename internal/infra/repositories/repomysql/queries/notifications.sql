-- name: UpsertNotification :exec
INSERT INTO notifications (id, alert_id, channel_id, message_id, status, checksum, sent_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    status     = VALUES(status),
    message_id = VALUES(message_id),
    sent_at    = VALUES(sent_at);

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = ? LIMIT 1;

-- name: GetNotificationByChecksum :one
SELECT * FROM notifications WHERE checksum = ? LIMIT 1;

-- name: UpdateNotificationStatus :exec
UPDATE notifications SET status = ? WHERE id = ?;

-- name: ListNotificationsByAlert :many
SELECT * FROM notifications WHERE alert_id = ? ORDER BY created_at DESC;
