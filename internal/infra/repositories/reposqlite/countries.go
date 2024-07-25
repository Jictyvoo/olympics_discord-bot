package reposqlite

import (
	"context"
	"database/sql"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) insertCountryCtx(
	ctx context.Context, dbQuery *dbgen.Queries, country entities.CountryInfo,
) (insertedID entities.Identifier, err error) {
	id, insertErr := dbQuery.InsertCountry(
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

	if insertErr != nil {
		return 0, insertErr
	}

	insertedID = entities.Identifier(id)
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
		if err != nil {
			_ = tx.Rollback()
		}
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
