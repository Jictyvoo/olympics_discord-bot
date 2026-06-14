-- name: UpsertResult :exec
INSERT INTO results (fixture_id, participant_id, position, score, raw_mark, outcome)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    position   = VALUES(position),
    score      = VALUES(score),
    raw_mark   = VALUES(raw_mark),
    outcome    = VALUES(outcome);

-- name: ListResultsByFixture :many
SELECT * FROM results WHERE fixture_id = ? ORDER BY position IS NULL, position ASC;

-- name: ListResultsByParticipant :many
SELECT * FROM results WHERE participant_id = ? ORDER BY created_at DESC;
