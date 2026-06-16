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

-- name: GetVenueByFixture :one
SELECT v.id, v.created_at, v.updated_at, v.provider_id, v.external_key, v.name, v.city, v.country_iso
FROM venues v
JOIN fixtures f ON f.venue_id = v.id
WHERE f.id = ? LIMIT 1;
