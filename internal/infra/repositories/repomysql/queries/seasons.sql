-- name: UpsertSeason :exec
INSERT INTO seasons (
    id, provider_id, external_key, name, starts_on, ends_on, competition_id
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name           = VALUES(name),
    starts_on      = VALUES(starts_on),
    ends_on        = VALUES(ends_on),
    competition_id = VALUES(competition_id);

-- name: GetSeason :one
SELECT * FROM seasons WHERE id = ? LIMIT 1;
