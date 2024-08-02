-- name: SaveResults :exec
INSERT INTO results (id, competitor_id, event_id, position, mark, medal_type, irm)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (competitor_id, event_id) DO UPDATE SET position   = excluded.position,
                                                    mark       = excluded.mark,
                                                    medal_type = excluded.medal_type,
                                                    irm        = excluded.irm;

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
