-- name: UpsertFixture :exec
INSERT INTO fixtures (
    id, provider_id, external_key, stage_id, group_id, venue_id,
    name, starts_at, ends_at, status, checksum
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name       = VALUES(name),
    starts_at  = VALUES(starts_at),
    ends_at    = VALUES(ends_at),
    status     = VALUES(status),
    checksum   = VALUES(checksum),
    group_id   = VALUES(group_id),
    venue_id   = VALUES(venue_id);

-- name: GetFixture :one
SELECT * FROM fixtures WHERE id = ? LIMIT 1;

-- name: GetFixtureByExternalKey :one
SELECT * FROM fixtures WHERE provider_id = ? AND external_key = ? LIMIT 1;

-- name: UpdateFixtureChecksum :exec
UPDATE fixtures SET checksum = ? WHERE id = ?;

-- name: ListFixturesByDay :many
SELECT * FROM fixtures
WHERE provider_id = ?
  AND starts_at >= ?
  AND starts_at < ?
ORDER BY starts_at ASC;

-- name: ListFixturesByStatus :many
SELECT * FROM fixtures WHERE status = ? ORDER BY starts_at ASC;

-- name: ListFixturesStartingBefore :many
SELECT * FROM fixtures
WHERE starts_at <= ?
  AND status NOT IN ('finished', 'cancelled')
ORDER BY starts_at ASC;

-- name: UpdateFixtureStatus :exec
UPDATE fixtures SET status = ? WHERE id = ?;
