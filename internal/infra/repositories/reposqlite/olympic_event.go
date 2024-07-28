package reposqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) upsertDiscipline(
	ctx context.Context, dbQuery *dbgen.Queries, name string,
) (disciplineID entities.Identifier, err error) {
	var foundID int64

	foundID, err = dbQuery.GetDisciplineIDByName(ctx, name)
	if foundID > 0 || !errors.Is(sql.ErrNoRows, err) {
		return entities.Identifier(foundID), nil
	}

	id, insertErr := dbQuery.InsertDiscipline(
		ctx, dbgen.InsertDisciplineParams{
			Name:        name,
			Description: nil,
		},
	)

	return entities.Identifier(id), insertErr
}

func (r RepoSQLite) SaveEvent(
	event entities.OlympicEvent, competitorIDs []entities.Identifier,
) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	dbQuery := r.queries

	// Insert and retrieve discipline ID
	disciplineID, err := r.upsertDiscipline(ctx, dbQuery, event.DisciplineName)
	if err != nil {
		return err
	}

	_, insertErr := dbQuery.SaveEvent(
		ctx, dbgen.SaveEventParams{
			EventName:    event.EventName,
			DisciplineID: int64(disciplineID),
			Phase:        event.Phase,
			Gender:       int64(event.Gender),
			StartAt:      event.StartAt,
			EndAt:        event.EndAt,
			Status:       string(event.Status),
		},
	)

	if insertErr != nil {
		return insertErr
	}

	// Create a result table row with competitors+event

	return nil
}
