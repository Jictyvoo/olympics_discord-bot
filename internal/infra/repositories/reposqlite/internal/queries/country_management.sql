-- name: InsertCountry :one
INSERT OR
REPLACE
INTO country_infos (created_at, updated_at, code, name, code_num, iso_code_len2, iso_code_len3,
                    ioc_code, population, area_km2, gdp_usd)
VALUES (datetime('now'), datetime('now'), ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: GetCountryByISOCode :one
SELECT id,
       created_at,
       updated_at,
       code,
       name,
       code_num,
       iso_code_len2,
       iso_code_len3,
       ioc_code,
       population,
       area_km2,
       gdp_usd
FROM country_infos
WHERE iso_code_len2 = ?
   OR iso_code_len3 = ?;
