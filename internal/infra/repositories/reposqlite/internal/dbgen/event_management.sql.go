// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: event_management.sql

package dbgen

import (
	"context"
	"time"
)

const GetDisciplineIDByName = `-- name: GetDisciplineIDByName :one
SELECT id
FROM olympic_disciplines
WHERE name = ?
`

func (q *Queries) GetDisciplineIDByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, GetDisciplineIDByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const GetEvent = `-- name: GetEvent :one
SELECT e.id                     as event_id,
       e.event_name,
       od.name                  as discipline_name,
       e.phase,
       e.gender,
       CAST(e.start_at AS TEXT) as start_at,
       CAST(e.end_at AS TEXT)   as end_at,
       e.status
FROM olympic_events e
         INNER JOIN
     olympic_disciplines od on e.discipline_id = od.id
WHERE e.id = ?
`

type GetEventRow struct {
	EventID        int64  `db:"event_id"`
	EventName      string `db:"event_name"`
	DisciplineName string `db:"discipline_name"`
	Phase          string `db:"phase"`
	Gender         int64  `db:"gender"`
	StartAt        string `db:"start_at"`
	EndAt          string `db:"end_at"`
	Status         string `db:"status"`
}

func (q *Queries) GetEvent(ctx context.Context, id int64) (GetEventRow, error) {
	row := q.db.QueryRowContext(ctx, GetEvent, id)
	var i GetEventRow
	err := row.Scan(
		&i.EventID,
		&i.EventName,
		&i.DisciplineName,
		&i.Phase,
		&i.Gender,
		&i.StartAt,
		&i.EndAt,
		&i.Status,
	)
	return i, err
}

const InsertDiscipline = `-- name: InsertDiscipline :one
INSERT INTO olympic_disciplines (name, description)
VALUES (?, ?)
ON CONFLICT DO NOTHING
RETURNING id
`

type InsertDisciplineParams struct {
	Name        string      `db:"name"`
	Description interface{} `db:"description"`
}

func (q *Queries) InsertDiscipline(ctx context.Context, arg InsertDisciplineParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, InsertDiscipline, arg.Name, arg.Description)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const LoadDayEvents = `-- name: LoadDayEvents :many
SELECT e.id                     as event_id,
       e.event_name,
       od.name                  as discipline_name,
       e.phase,
       e.gender,
       CAST(e.start_at AS TEXT) as start_at,
       CAST(e.end_at AS TEXT)   as end_at,
       e.status
FROM olympic_events e
         INNER JOIN
     olympic_disciplines od on e.discipline_id = od.id
WHERE e.start_at >= ?
  AND e.end_at <= ?
ORDER BY e.start_at
`

type LoadDayEventsParams struct {
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
}

type LoadDayEventsRow struct {
	EventID        int64  `db:"event_id"`
	EventName      string `db:"event_name"`
	DisciplineName string `db:"discipline_name"`
	Phase          string `db:"phase"`
	Gender         int64  `db:"gender"`
	StartAt        string `db:"start_at"`
	EndAt          string `db:"end_at"`
	Status         string `db:"status"`
}

func (q *Queries) LoadDayEvents(ctx context.Context, arg LoadDayEventsParams) ([]LoadDayEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, LoadDayEvents, arg.StartAt, arg.EndAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LoadDayEventsRow{}
	for rows.Next() {
		var i LoadDayEventsRow
		if err := rows.Scan(
			&i.EventID,
			&i.EventName,
			&i.DisciplineName,
			&i.Phase,
			&i.Gender,
			&i.StartAt,
			&i.EndAt,
			&i.Status,
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

const SaveEvent = `-- name: SaveEvent :one
INSERT INTO olympic_events (event_name, discipline_id, phase, gender, start_at, end_at, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (event_name, discipline_id, phase, gender) DO UPDATE SET status=excluded.status,
                                                                     start_at=excluded.start_at,
                                                                     end_at=excluded.end_at
RETURNING id
`

type SaveEventParams struct {
	EventName    string    `db:"event_name"`
	DisciplineID int64     `db:"discipline_id"`
	Phase        string    `db:"phase"`
	Gender       int64     `db:"gender"`
	StartAt      time.Time `db:"start_at"`
	EndAt        time.Time `db:"end_at"`
	Status       string    `db:"status"`
}

func (q *Queries) SaveEvent(ctx context.Context, arg SaveEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, SaveEvent,
		arg.EventName,
		arg.DisciplineID,
		arg.Phase,
		arg.Gender,
		arg.StartAt,
		arg.EndAt,
		arg.Status,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
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
