package reposqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) upsertDiscipline(
	ctx context.Context, dbQuery *dbgen.Queries, name string,
) (disciplineID entities.Identifier, err error) {
	var foundID int64

	foundID, err = dbQuery.GetDisciplineIDByName(ctx, name)
	if foundID > 0 || (err != nil && !errors.Is(sql.ErrNoRows, err)) {
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

	eventID, insertErr := dbQuery.SaveEvent(
		ctx, dbgen.SaveEventParams{
			EventName:    event.EventName,
			DisciplineID: int64(disciplineID),
			Phase:        event.Phase,
			Gender:       int64(event.Gender),
			SessionCode:  event.SessionCode,
			StartAt:      event.StartAt,
			EndAt:        event.EndAt,
			Status:       string(event.Status),
		},
	)

	if insertErr != nil || eventID == 0 {
		return insertErr
	}

	// Create a result table row with competitors+event
	tx, txErr := r.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	dbQuery = dbQuery.WithTx(tx)
	for _, competitorID := range competitorIDs {
		err = dbQuery.SaveResults(
			ctx, dbgen.SaveResultsParams{
				ID:           uuid.New().String(),
				CompetitorID: int64(competitorID),
				EventID:      eventID,
				Position:     nil,
				Mark:         nil,
				MedalType:    nil,
				Irm:          "",
			},
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func parseStartEndTimes(resultEvent *entities.OlympicEvent, startAt, endAt string) {
	const layout = "2006-01-02 15:04:05 -0700 -0700"
	resultEvent.StartAt, _ = time.Parse(layout, startAt)
	resultEvent.EndAt, _ = time.Parse(layout, endAt)
}

func (r RepoSQLite) LoadDayEvents(from time.Time) ([]entities.OlympicEvent, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbQuery := r.queries
	foundEvents, err := dbQuery.LoadDayEvents(
		ctx, dbgen.LoadDayEventsParams{
			StartAt: from,
			EndAt: time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC).
				Add(24 * time.Hour),
		},
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	eventList := make([]entities.OlympicEvent, len(foundEvents))
	for index, foundEvent := range foundEvents {
		eventID := entities.Identifier(foundEvent.EventID)
		competitors, searchErr := r.loadEventCompetitorsCtx(
			ctx, dbQuery, eventID,
		)
		if searchErr != nil {
			return nil, searchErr
		}

		eventList[index] = entities.OlympicEvent{
			ID:             eventID,
			EventName:      foundEvent.EventName,
			DisciplineName: foundEvent.DisciplineName,
			Phase:          foundEvent.Phase,
			Gender:         entities.Gender(foundEvent.Gender),
			SessionCode:    foundEvent.SessionCode,
			Status:         entities.EventStatus(foundEvent.Status),
			Competitors:    competitors,
		}

		parseStartEndTimes(&eventList[index], foundEvent.StartAt, foundEvent.EndAt)
	}

	return eventList, nil
}

func (r RepoSQLite) LoadEvent(id entities.Identifier) (entities.OlympicEvent, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbQuery := r.queries
	foundEvent, err := dbQuery.GetEvent(
		ctx, int64(id),
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entities.OlympicEvent{}, err
	}

	eventID := entities.Identifier(foundEvent.EventID)
	competitors, searchErr := r.loadEventCompetitorsCtx(
		ctx, dbQuery, eventID,
	)
	if searchErr != nil {
		return entities.OlympicEvent{}, searchErr
	}

	resultEvent := entities.OlympicEvent{
		ID:             eventID,
		EventName:      foundEvent.EventName,
		DisciplineName: foundEvent.DisciplineName,
		Phase:          foundEvent.Phase,
		Gender:         entities.Gender(foundEvent.Gender),
		SessionCode:    foundEvent.SessionCode,
		Status:         entities.EventStatus(foundEvent.Status),
		Competitors:    competitors,
	}

	parseStartEndTimes(&resultEvent, foundEvent.StartAt, foundEvent.EndAt)

	return resultEvent, nil
}
