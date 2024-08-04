-- name: GetDisciplineIDByName :one
SELECT id, code
FROM olympic_disciplines
WHERE name = ?;


-- name: InsertDiscipline :one
INSERT INTO olympic_disciplines (name, description, code)
VALUES (?, ?, ?)
ON CONFLICT(name) DO UPDATE SET code=excluded.code
RETURNING id;


-- name: SaveEvent :one
INSERT INTO olympic_events (event_name, discipline_id, phase, gender, session_code,
                            start_at, end_at, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT (event_name, discipline_id, phase, gender, session_code) DO UPDATE SET status=excluded.status,
                                                                                   start_at=excluded.start_at,
                                                                                   end_at=excluded.end_at,
                                                                                   updated_at=excluded.updated_at
RETURNING id;


-- name: LoadDayEvents :many
SELECT e.id                     as event_id,
       e.event_name,
       od.name                  as discipline_name,
       od.code                  as discipline_code,
       e.phase,
       e.gender,
       e.session_code,
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
       od.code                  as discipline_code,
       e.phase,
       e.gender,
       e.session_code,
       CAST(e.start_at AS TEXT) as start_at,
       CAST(e.end_at AS TEXT)   as end_at,
       e.status
FROM olympic_events e
         INNER JOIN olympic_disciplines od on e.discipline_id = od.id
WHERE e.id = ?;
