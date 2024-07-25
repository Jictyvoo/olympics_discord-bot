package reposqlite

import (
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
	id, insertErr := dbQuery.SaveCompetitor(
		ctx, dbgen.SaveCompetitorParams{
			Code:      competitor.Code,
			Name:      competitor.Name,
			CountryID: int64(countryID),
		},
	)

	if insertErr != nil {
		return 0, insertErr
	}

	insertedID := entities.Identifier(id)
	return insertedID, nil
}
