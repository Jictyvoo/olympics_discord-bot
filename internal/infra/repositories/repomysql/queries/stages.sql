-- name: UpsertStage :exec
INSERT INTO stages (
    id, provider_id, external_key, name, ord, season_id
) VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name      = VALUES(name),
    ord       = VALUES(ord),
    season_id = VALUES(season_id);

-- name: GetStage :one
SELECT * FROM stages WHERE id = ? LIMIT 1;
