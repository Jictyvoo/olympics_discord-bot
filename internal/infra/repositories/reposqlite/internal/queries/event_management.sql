-- name: GetDisciplineIDByName :one
SELECT id
FROM olympic_disciplines
WHERE name = ?;


-- name: InsertDiscipline :one
INSERT INTO olympic_disciplines (name, description)
VALUES (?, ?)
ON CONFLICT DO NOTHING
RETURNING id;


-- name: SaveEvent :one
INSERT INTO olympic_events (event_name, discipline_id, phase, gender, start_at, end_at, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (event_name, discipline_id, phase, gender) DO UPDATE SET status=excluded.status,
                                                                     start_at=excluded.start_at,
                                                                     end_at=excluded.end_at
RETURNING id;


-- name: SaveResults :exec
INSERT INTO results (id, competitor_id, event_id, position, mark, medal_type, irm)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (competitor_id, event_id) DO UPDATE SET position   = excluded.position,
                                                    mark       = excluded.mark,
                                                    medal_type = excluded.medal_type,
                                                    irm        = excluded.irm;

-- name: LoadDayEvents :many
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
ORDER BY e.start_at;


-- name: GetEvent :one
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
WHERE e.id = ?;
