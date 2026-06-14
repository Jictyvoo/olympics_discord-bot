-- name: UpsertStanding :exec
INSERT INTO standings (
    stage_id, participant_id, `rank`, points
) VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    `rank` = VALUES(`rank`),
    points = VALUES(points);

-- name: ListStandingsByStage :many
SELECT * FROM standings WHERE stage_id = ? ORDER BY `rank` ASC;
