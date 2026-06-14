-- name: UpsertGroup :exec
INSERT INTO groups (
    id, provider_id, external_key, name, stage_id,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    name       = excluded.name,
    stage_id   = excluded.stage_id,
    updated_at = DATETIME('now');

-- name: GetGroup :one
SELECT * FROM groups WHERE id = ? LIMIT 1;
