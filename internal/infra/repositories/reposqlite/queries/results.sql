-- name: UpsertResult :exec
INSERT INTO results (fixture_id, participant_id, position, score, raw_mark, outcome, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(fixture_id, participant_id) DO UPDATE SET
    position   = excluded.position,
    score      = excluded.score,
    raw_mark   = excluded.raw_mark,
    outcome    = excluded.outcome,
    updated_at = DATETIME('now');
