-- name: SaveCompetitor :one
INSERT OR
REPLACE INTO competitors (code, name, country_id)
VALUES (?, ?, ?)
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


-- name: GetCompetitorByName :one
SELECT c.id,
       c.code,
       c.name,
       c.country_id
FROM competitors c
WHERE c.id = ?;
