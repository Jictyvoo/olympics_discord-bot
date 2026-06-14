-- name: UpsertStage :exec
INSERT INTO stages (
    id, provider_id, external_key, name, ord, season_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    name       = excluded.name,
    ord        = excluded.ord,
    season_id  = excluded.season_id,
    updated_at = DATETIME('now');

-- name: GetStage :one
SELECT * FROM stages WHERE id = ? LIMIT 1;
