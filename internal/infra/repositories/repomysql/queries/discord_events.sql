-- name: UpsertDiscordEvent :exec
INSERT INTO discord_events (fixture_id, guild_id, discord_event_id, status, last_checksum)
VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    discord_event_id = VALUES(discord_event_id),
    status           = VALUES(status),
    last_checksum    = VALUES(last_checksum);

-- name: GetDiscordEventByFixture :one
SELECT * FROM discord_events WHERE fixture_id = ? AND guild_id = ? LIMIT 1;

-- name: GetDiscordEventByDiscordID :one
SELECT * FROM discord_events WHERE discord_event_id = ? LIMIT 1;

-- name: UpdateDiscordEventStatus :exec
UPDATE discord_events SET status = ? WHERE fixture_id = ? AND guild_id = ?;

-- name: ListActiveDiscordEvents :many
SELECT * FROM discord_events
WHERE guild_id = ? AND status NOT IN ('completed', 'cancelled')
ORDER BY updated_at DESC;
