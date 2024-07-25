-- name: SaveEvent :one
INSERT INTO olympic_events (event_name, discipline_name, phase, gender, start_at, end_at, status)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id;
