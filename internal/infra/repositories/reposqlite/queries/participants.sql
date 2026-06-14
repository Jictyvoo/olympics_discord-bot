-- name: SaveCompetitor :one
INSERT INTO competitors (code, name, country_id)
VALUES (?, ?, ?)
ON CONFLICT (code, name, country_id) DO NOTHING
RETURNING id;


-- name: GetCompetitorByCountry :many
SELECT c.id,
       c.code,
       c.name,
       c.country_id
FROM competitors c
         JOIN
     country_infos ci ON c.country_id = ci.id
WHERE ci.iso_code_len2 = ?
   OR ci.iso_code_len3 = ?;


-- name: GetCompetitor :one
SELECT c.id,
       c.code,
       c.name,
       c.country_id
FROM competitors c
WHERE c.code = ?
  AND c.name = ?
  AND c.country_id = ?;


-- name: GetEventCompetitors :many
SELECT c.name,
       c.code,
       ci.name AS country_name,
       ci.code AS country_code,
       ci.ioc_code,
       ci.iso_code_len2,
       ci.iso_code_len3
FROM competitors c
         INNER JOIN
     results r ON c.id = r.competitor_id
         INNER JOIN main.country_infos ci on ci.id = c.country_id
WHERE r.event_id = ?;
