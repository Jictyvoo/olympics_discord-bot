-- name: UpsertStanding :exec
INSERT INTO standings (
    stage_id, participant_id, rank, points, created_at, updated_at
) VALUES (?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(stage_id, participant_id) DO UPDATE SET
    rank       = excluded.rank,
    points     = excluded.points,
    updated_at = DATETIME('now');

-- name: ListStandingsByStage :many
SELECT * FROM standings WHERE stage_id = ? ORDER BY rank ASC;
