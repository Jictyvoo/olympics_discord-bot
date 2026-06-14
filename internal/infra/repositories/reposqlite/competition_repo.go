package reposqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

func (r RepoSQLite) upsertDiscipline(
	ctx context.Context, dbQuery *dbgen.Queries,
	newDisc entities.Discipline,
) (disciplineID entities.Identifier, err error) {
	var foundDiscipline dbgen.GetDisciplineIDByNameRow

	foundDiscipline, err = dbQuery.GetDisciplineIDByName(ctx, newDisc.Name)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	disciplineID = entities.Identifier(foundDiscipline.ID)
	if (foundDiscipline.Code == "" && newDisc.Code != "") || foundDiscipline.ID <= 0 {
		id, insertErr := dbQuery.InsertDiscipline(
			ctx, dbgen.InsertDisciplineParams{
				Name:        newDisc.Name,
				Code:        newDisc.Code,
				Description: newDisc.Description,
			},
		)
		disciplineID = entities.Identifier(id)
		err = insertErr
	}

	return
}

func (r RepoSQLite) SaveDisciplines(disciplineList []entities.Discipline) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Prepare the statement
	tx, txErr := r.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	dbQuery := r.queries.WithTx(tx)

	for _, discipline := range disciplineList {
		disciplineID, err := r.upsertDiscipline(ctx, dbQuery, discipline)
		if err != nil {
			return err
		}

		if disciplineID == 0 {
			slog.Warn(
				"Recover a zero discipline ID during insert",
				slog.Any("discipline", discipline),
			)
		}
	}

	return tx.Commit()
}
