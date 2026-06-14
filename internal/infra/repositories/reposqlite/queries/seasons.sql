-- name: UpsertSeason :exec
INSERT INTO seasons (
    id, provider_id, external_key, name, starts_on, ends_on, competition_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    name           = excluded.name,
    starts_on      = excluded.starts_on,
    ends_on        = excluded.ends_on,
    competition_id = excluded.competition_id,
    updated_at     = DATETIME('now');

-- name: GetSeason :one
SELECT * FROM seasons WHERE id = ? LIMIT 1;
