-- name: InsertCountry :one
INSERT INTO country_infos (code, name, code_num, iso_code_len2, iso_code_len3, ioc_code,
                           population, area_km2, gdp_usd, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(ioc_code) DO NOTHING
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
