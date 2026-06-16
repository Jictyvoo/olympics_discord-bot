-- name: UpsertCompetition :exec
INSERT INTO competitions (
    id, provider_id, external_key, code, name, discipline
) VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    code       = VALUES(code),
    name       = VALUES(name),
    discipline = VALUES(discipline);

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
LEFT JOIN `groups` g  ON g.id = f.group_id
WHERE f.id = ? LIMIT 1;
