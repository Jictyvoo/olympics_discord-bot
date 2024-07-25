package reposqlite

import (
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) SaveEvent(
	event entities.OlympicEvent, competitorIDs []entities.Identifier,
) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	dbQuery := r.queries
	_, insertErr := dbQuery.SaveEvent(
		ctx, dbgen.SaveEventParams{
			EventName:      event.EventName,
			DisciplineName: event.DisciplineName,
			Phase:          event.Phase,
			Gender:         int64(event.Gender),
			StartAt:        event.StartAt,
			EndAt:          event.EndAt,
			Status:         string(event.Status),
		},
	)

	if insertErr != nil {
		return insertErr
	}

	// Create a result table row with competitors+event

	return nil
}
