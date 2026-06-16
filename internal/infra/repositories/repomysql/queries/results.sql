-- name: UpsertResult :exec
INSERT INTO results (fixture_id, participant_id, position, score, raw_mark, outcome)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    position   = VALUES(position),
    score      = VALUES(score),
    raw_mark   = VALUES(raw_mark),
    outcome    = VALUES(outcome);
