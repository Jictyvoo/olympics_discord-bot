// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: event_results.sql

package dbgen

import (
	"context"
)

const GetEventResults = `-- name: GetEventResults :many
SELECT r.event_id      AS event_id,
       r.competitor_id AS competitor_id,
       c.code          AS competitor_code,
       r.position,
       r.mark,
       r.medal_type,
       r.irm
FROM olympic_events e
         INNER JOIN results r ON e.id = r.event_id
         INNER JOIN competitors c on r.competitor_id = c.id
WHERE e.id = ?
`

type GetEventResultsRow struct {
	EventID        int64       `db:"event_id"`
	CompetitorID   int64       `db:"competitor_id"`
	CompetitorCode string      `db:"competitor_code"`
	Position       interface{} `db:"position"`
	Mark           interface{} `db:"mark"`
	MedalType      interface{} `db:"medal_type"`
	Irm            string      `db:"irm"`
}

func (q *Queries) GetEventResults(ctx context.Context, id int64) ([]GetEventResultsRow, error) {
	rows, err := q.db.QueryContext(ctx, GetEventResults, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEventResultsRow{}
	for rows.Next() {
		var i GetEventResultsRow
		if err := rows.Scan(
			&i.EventID,
			&i.CompetitorID,
			&i.CompetitorCode,
			&i.Position,
			&i.Mark,
			&i.MedalType,
			&i.Irm,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const SaveResults = `-- name: SaveResults :exec
INSERT INTO results (id, competitor_id, event_id, position, mark, medal_type, irm)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (competitor_id, event_id) DO UPDATE SET position   = excluded.position,
                                                    mark       = excluded.mark,
                                                    medal_type = excluded.medal_type,
                                                    irm        = excluded.irm
`

type SaveResultsParams struct {
	ID           interface{} `db:"id"`
	CompetitorID int64       `db:"competitor_id"`
	EventID      int64       `db:"event_id"`
	Position     interface{} `db:"position"`
	Mark         interface{} `db:"mark"`
	MedalType    interface{} `db:"medal_type"`
	Irm          string      `db:"irm"`
}

func (q *Queries) SaveResults(ctx context.Context, arg SaveResultsParams) error {
	_, err := q.db.ExecContext(ctx, SaveResults,
		arg.ID,
		arg.CompetitorID,
		arg.EventID,
		arg.Position,
		arg.Mark,
		arg.MedalType,
		arg.Irm,
	)
	return err
}
