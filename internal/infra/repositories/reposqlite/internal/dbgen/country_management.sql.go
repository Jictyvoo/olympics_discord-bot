// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: country_management.sql

package dbgen

import (
	"context"
)

const GetCountryByISOCode = `-- name: GetCountryByISOCode :one
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
   OR iso_code_len3 = ?
`

type GetCountryByISOCodeParams struct {
	IsoCodeLen2 interface{} `db:"iso_code_len2"`
	IsoCodeLen3 string      `db:"iso_code_len3"`
}

func (q *Queries) GetCountryByISOCode(ctx context.Context, arg GetCountryByISOCodeParams) (CountryInfo, error) {
	row := q.db.QueryRowContext(ctx, GetCountryByISOCode, arg.IsoCodeLen2, arg.IsoCodeLen3)
	var i CountryInfo
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Code,
		&i.Name,
		&i.CodeNum,
		&i.IsoCodeLen2,
		&i.IsoCodeLen3,
		&i.IocCode,
		&i.Population,
		&i.AreaKm2,
		&i.GdpUsd,
	)
	return i, err
}

const InsertCountry = `-- name: InsertCountry :one
INSERT OR
REPLACE
INTO country_infos (created_at, updated_at, code, name, code_num, iso_code_len2, iso_code_len3,
                    ioc_code, population, area_km2, gdp_usd)
VALUES (datetime('now'), datetime('now'), ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id
`

type InsertCountryParams struct {
	Code        string      `db:"code"`
	Name        string      `db:"name"`
	CodeNum     string      `db:"code_num"`
	IsoCodeLen2 interface{} `db:"iso_code_len2"`
	IsoCodeLen3 string      `db:"iso_code_len3"`
	IocCode     string      `db:"ioc_code"`
	Population  interface{} `db:"population"`
	AreaKm2     interface{} `db:"area_km2"`
	GdpUsd      interface{} `db:"gdp_usd"`
}

func (q *Queries) InsertCountry(ctx context.Context, arg InsertCountryParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, InsertCountry,
		arg.Code,
		arg.Name,
		arg.CodeNum,
		arg.IsoCodeLen2,
		arg.IsoCodeLen3,
		arg.IocCode,
		arg.Population,
		arg.AreaKm2,
		arg.GdpUsd,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}
