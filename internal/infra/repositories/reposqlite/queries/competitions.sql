-- name: UpsertCompetition :exec
INSERT INTO competitions (
    id, provider_id, external_key, code, name, discipline,
    created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, DATETIME('now'), DATETIME('now'))
ON CONFLICT(id) DO UPDATE SET
    code       = excluded.code,
    name       = excluded.name,
    discipline = excluded.discipline,
    updated_at = DATETIME('now');

-- name: GetCompetition :one
SELECT * FROM competitions WHERE id = ? LIMIT 1;

-- name: GetFixtureContext :one
SELECT
    c.code        AS competition_code,
    c.name        AS competition_name,
    c.discipline  AS discipline,
    st.name       AS stage_name,
    st.ord        AS stage_ord,
    g.name        AS group_name
FROM fixtures f
JOIN stages st      ON st.id = f.stage_id
JOIN seasons s      ON s.id = st.season_id
JOIN competitions c ON c.id = s.competition_id
LEFT JOIN groups g  ON g.id = f.group_id
WHERE f.id = ? LIMIT 1;
