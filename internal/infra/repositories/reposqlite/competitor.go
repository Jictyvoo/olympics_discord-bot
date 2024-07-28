package reposqlite

import (
	"database/sql"
	"errors"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) SaveCompetitor(
	competitor entities.OlympicCompetitors,
) (entities.Identifier, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	dbQuery := r.queries

	countryID, err := r.insertCountryCtx(ctx, dbQuery, competitor.CountryInfo)
	if err != nil {
		return 0, err
	}

	var id int64
	// Try to find competitor
	foundCompetitor, foundErr := dbQuery.GetCompetitor(
		ctx, dbgen.GetCompetitorParams{
			Code:      competitor.Code,
			Name:      competitor.Name,
			CountryID: int64(countryID),
		},
	)
	if !errors.Is(foundErr, sql.ErrNoRows) || foundCompetitor.ID > 0 {
		return entities.Identifier(foundCompetitor.ID), foundErr
	}

	id, err = dbQuery.SaveCompetitor(
		ctx, dbgen.SaveCompetitorParams{
			Code:      competitor.Code,
			Name:      competitor.Name,
			CountryID: int64(countryID),
		},
	)

	insertedID := entities.Identifier(id)
	return insertedID, err
}
