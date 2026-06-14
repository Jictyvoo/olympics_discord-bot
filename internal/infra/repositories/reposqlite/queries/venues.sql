-- name: UpsertVenue :exec
INSERT INTO venues (
    id, provider_id, external_key, name, city, country_iso,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    name        = excluded.name,
    city        = excluded.city,
    country_iso = excluded.country_iso,
    updated_at  = DATETIME('now');

-- name: GetVenue :one
SELECT * FROM venues WHERE id = ? LIMIT 1;
