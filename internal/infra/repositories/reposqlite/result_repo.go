package reposqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) saveResultsCtx(
	ctx context.Context, dbQuery *dbgen.Queries,
	competitorResultsByIDs map[entities.Identifier]*entities.Results, eventID int64,
) (err error) {
	tx, txErr := r.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	competitorIDs := make([]int64, 0, len(competitorResultsByIDs))
	dbQuery = dbQuery.WithTx(tx)
	for competitorID, result := range competitorResultsByIDs {
		competitorIDs = append(competitorIDs, int64(competitorID))
		saveParams := dbgen.SaveResultsParams{
			ID:           uuid.New().String(),
			CompetitorID: int64(competitorID),
			EventID:      eventID,
			Position:     nil,
			Mark:         nil,
			MedalType:    nil,
			Irm:          "",
		}
		if result != nil {
			saveParams.Position = result.Position
			saveParams.Mark = result.Mark
			saveParams.MedalType = result.MedalType
			saveParams.Irm = result.Irm
		}

		err = dbQuery.SaveResults(ctx, saveParams)
		if err != nil {
			return err
		}
	}

	err = dbQuery.DeleteResultsWithCompetitors(
		ctx, dbgen.DeleteResultsWithCompetitorsParams{
			EventID: eventID, CompetitorIds: competitorIDs,
		},
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r RepoSQLite) loadResults(
	ctx context.Context, eventID entities.Identifier,
) (resultPerCompetitor map[string]entities.Results, err error) {
	dbQuery := r.queries
	var results []dbgen.GetEventResultsRow
	if results, err = dbQuery.GetEventResults(ctx, int64(eventID)); errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if err != nil {
		return map[string]entities.Results{}, err
	}

	resultPerCompetitor = make(map[string]entities.Results, len(results))

	for _, result := range results {
		loadResult := entities.Results{
			Irm: result.Irm,
		}
		loadResult.Position, _ = result.Position.(string)
		loadResult.Mark, _ = result.Mark.(string)
		medalType, _ := result.MedalType.(string)
		loadResult.MedalType = entities.Medal(medalType)

		if loadResult.MedalType != entities.MedalNoMedal || loadResult.Mark != "" {
			resultPerCompetitor[result.CompetitorCode] = loadResult
		}
	}

	return
}
