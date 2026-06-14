-- name: UpsertCountry :exec
INSERT IGNORE INTO countries (iso2, iso3, ioc_code, name, code_num, population, area_km2, gdp_usd)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCountryByISO2 :one
SELECT * FROM countries WHERE iso2 = ? LIMIT 1;

-- name: GetCountryByISO3 :one
SELECT * FROM countries WHERE iso3 = ? LIMIT 1;

-- name: GetCountryByIOC :one
SELECT * FROM countries WHERE ioc_code = ? LIMIT 1;

-- name: ListCountries :many
SELECT * FROM countries ORDER BY name ASC;
