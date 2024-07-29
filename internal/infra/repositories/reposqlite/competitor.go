package reposqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) loadEventCompetitorsCtx(
	ctx context.Context, dbQuery *dbgen.Queries, eventID entities.Identifier,
) ([]entities.OlympicCompetitors, error) {
	foundCompetitors, err := dbQuery.GetEventCompetitors(ctx, int64(eventID))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	competitorList := make([]entities.OlympicCompetitors, len(foundCompetitors))
	for index, competitor := range foundCompetitors {
		isoCode2, _ := competitor.IsoCodeLen2.(string)
		competitorList[index] = entities.OlympicCompetitors{
			Code:        competitor.Code,
			CountryCode: competitor.CountryCode,
			Name:        competitor.Name,
			Country: entities.CountryInfo{
				Name:    competitor.CountryName,
				ISOCode: [2]string{isoCode2, competitor.IsoCodeLen3},
				IOCCode: competitor.IocCode,
			},
		}
	}

	return competitorList, nil
}

func (r RepoSQLite) SaveCompetitor(
	competitor entities.OlympicCompetitors,
) (entities.Identifier, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	dbQuery := r.queries

	countryID, err := r.insertCountryCtx(ctx, dbQuery, competitor.Country)
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
	if (foundErr != nil && !errors.Is(foundErr, sql.ErrNoRows)) || foundCompetitor.ID > 0 {
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
