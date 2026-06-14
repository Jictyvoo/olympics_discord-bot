-- name: UpsertGroup :exec
INSERT INTO `groups` (
    id, provider_id, external_key, name, stage_id
) VALUES (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name     = VALUES(name),
    stage_id = VALUES(stage_id);

-- name: GetGroup :one
SELECT * FROM `groups` WHERE id = ? LIMIT 1;
