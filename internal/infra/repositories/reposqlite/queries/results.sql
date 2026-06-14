-- name: SaveResults :exec
INSERT INTO results (id, competitor_id, event_id, position, mark, medal_type, irm, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT (competitor_id, event_id) DO UPDATE SET position   = excluded.position,
                                                    mark       = excluded.mark,
                                                    medal_type = excluded.medal_type,
                                                    irm        = excluded.irm,
                                                    updated_at = excluded.updated_at;

-- name: GetEventResults :many
SELECT r.event_id      AS event_id,
       r.competitor_id AS competitor_id,
       c.code          AS competitor_code,
       r.position,
       r.mark,
       r.medal_type,
       r.irm
FROM olympic_events e
         INNER JOIN results r ON e.id = r.event_id
         INNER JOIN competitors c on r.competitor_id = c.id
WHERE e.id = ?
ORDER BY r.mark DESC, r.medal_type;


-- name: DeleteResultsWithCompetitors :exec
-- noinspection SqlResolve @ any/"sqlc"
DELETE
FROM results
WHERE results.event_id = ?
  AND results.competitor_id NOT IN (sqlc.slice('competitor_ids'))
