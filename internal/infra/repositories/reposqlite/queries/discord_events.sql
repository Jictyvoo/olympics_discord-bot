-- name: UpsertDiscordEvent :exec
INSERT INTO discord_events (fixture_id, guild_id, discord_event_id, status, last_checksum, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(fixture_id, guild_id) DO UPDATE SET
    discord_event_id = excluded.discord_event_id,
    status           = excluded.status,
    last_checksum    = excluded.last_checksum,
    updated_at       = DATETIME('now');

-- name: GetDiscordEventByFixture :one
SELECT * FROM discord_events WHERE fixture_id = ? AND guild_id = ? LIMIT 1;

-- name: GetDiscordEventByDiscordID :one
SELECT * FROM discord_events WHERE discord_event_id = ? LIMIT 1;

-- name: UpdateDiscordEventStatus :exec
UPDATE discord_events SET status = ?, updated_at = DATETIME('now')
WHERE fixture_id = ? AND guild_id = ?;

-- name: ListActiveDiscordEvents :many
SELECT * FROM discord_events
WHERE guild_id = ? AND status NOT IN ('completed', 'cancelled')
ORDER BY updated_at DESC;
