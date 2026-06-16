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

-- name: GetVenueByFixture :one
SELECT v.id, v.created_at, v.updated_at, v.provider_id, v.external_key, v.name, v.city, v.country_iso
FROM venues v
JOIN fixtures f ON f.venue_id = v.id
WHERE f.id = ? LIMIT 1;
