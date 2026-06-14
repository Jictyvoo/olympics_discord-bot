-- name: AddSubscription :exec
INSERT INTO subscriptions (guild_id, user_id, kind, value, created_at, updated_at)
VALUES (?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(guild_id, user_id, kind, value) DO UPDATE SET
    updated_at = DATETIME('now');

-- name: RemoveSubscription :exec
DELETE FROM subscriptions
WHERE guild_id = ? AND user_id = ? AND kind = ? AND value = ?;

-- name: ListSubscriptionsByGuild :many
SELECT * FROM subscriptions WHERE guild_id = ? ORDER BY id ASC;

-- name: ListSubscriptionsByGuildUser :many
SELECT * FROM subscriptions WHERE guild_id = ? AND user_id = ? ORDER BY id ASC;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions ORDER BY id ASC;
