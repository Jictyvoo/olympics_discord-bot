-- name: UpsertVenue :exec
INSERT INTO venues (
    id, provider_id, external_key, name, city, country_iso
) VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name        = VALUES(name),
    city        = VALUES(city),
    country_iso = VALUES(country_iso);

-- name: GetVenue :one
SELECT * FROM venues WHERE id = ? LIMIT 1;
