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
