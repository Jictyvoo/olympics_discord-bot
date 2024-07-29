package reposqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) insertCountryCtx(
	ctx context.Context, dbQuery *dbgen.Queries, country entities.CountryInfo,
) (countryID entities.Identifier, err error) {
	var foundCountry dbgen.CountryInfo
	foundCountry, err = dbQuery.GetCountryByISOCode(
		ctx, dbgen.GetCountryByISOCodeParams{
			IsoCodeLen2: country.ISOCode[0],
			IsoCodeLen3: country.ISOCode[1],
		},
	)

	if (err != nil && !errors.Is(err, sql.ErrNoRows)) || foundCountry.ID > 0 {
		countryID = entities.Identifier(foundCountry.ID)
		if err != nil {
			return
		}
		return
	}

	var insertedID int64
	insertedID, err = dbQuery.InsertCountry(
		ctx, dbgen.InsertCountryParams{
			Code:        country.CodeNum,
			Name:        country.Name,
			CodeNum:     country.CodeNum,
			IsoCodeLen2: country.ISOCode[0],
			IsoCodeLen3: country.ISOCode[1],
			IocCode:     country.IOCCode,
			Population:  country.Population,
			AreaKm2:     country.AreaKm2,
			GdpUsd:      country.GDPUSD,
		},
	)

	countryID = entities.Identifier(insertedID)
	if err != nil {
		println(err.Error())
	}
	return
}

func (r RepoSQLite) InsertCountry(
	country entities.CountryInfo,
) (insertedID entities.Identifier, err error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	dbQuery := r.queries
	return r.insertCountryCtx(ctx, dbQuery, country)
}

func (r RepoSQLite) InsertCountries(
	countries []entities.CountryInfo,
) (insertedIDs []entities.Identifier, err error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Initialize a slice to hold the returned identifiers
	insertedIDs = make([]entities.Identifier, len(countries))

	// Start a transaction
	var tx *sql.Tx
	if tx, err = r.conn.BeginTx(ctx, nil); err != nil {
		return
	}

	// Rollback the transaction in case of error
	defer func() {
		_ = tx.Rollback()
	}()

	// Prepare the statement
	dbQuery := r.queries.WithTx(tx)
	for index, country := range countries {
		insertedIDs[index], err = r.insertCountryCtx(ctx, dbQuery, country)
		if err != nil {
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	return
}
