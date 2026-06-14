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

-- name: GetCompetitionByFixture :one
SELECT c.id, c.created_at, c.updated_at, c.provider_id, c.external_key, c.code, c.name, c.discipline
FROM competitions c
JOIN seasons s ON s.competition_id = c.id
JOIN stages st ON st.season_id = s.id
JOIN fixtures f ON f.stage_id = st.id
WHERE f.id = ? LIMIT 1;
