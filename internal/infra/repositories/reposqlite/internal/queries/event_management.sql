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
INSERT INTO results (competitor_id, event_id, position, mark, medal_type, irm)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (competitor_id, event_id) DO UPDATE SET position   = excluded.position,
                                                    mark       = excluded.mark,
                                                    medal_type = excluded.medal_type,
                                                    irm        = excluded.irm;
