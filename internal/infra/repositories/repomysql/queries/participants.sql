-- name: UpsertParticipant :exec
INSERT INTO participants (
    id, provider_id, external_key, kind, name, code, country_iso, gender
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    name        = VALUES(name),
    code        = VALUES(code),
    country_iso = VALUES(country_iso),
    gender      = VALUES(gender);

-- name: GetParticipant :one
SELECT * FROM participants WHERE id = ? LIMIT 1;

-- name: GetParticipantByExternalKey :one
SELECT * FROM participants WHERE provider_id = ? AND external_key = ? LIMIT 1;

-- name: UpsertFixtureParticipant :exec
INSERT INTO fixture_participants (fixture_id, participant_id, role)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE role = VALUES(role);

-- name: ListFixtureParticipants :many
SELECT fp.fixture_id, fp.participant_id, fp.role
FROM fixture_participants fp
WHERE fp.fixture_id = ?;

-- name: ListFixtureCompetitors :many
SELECT
    p.id, p.provider_id, p.external_key, p.kind, p.name, p.code, p.country_iso, p.gender,
    fp.role,
    COALESCE((
        SELECT ct.iso2 FROM countries ct
        WHERE ct.iso3 = p.country_iso OR ct.ioc_code = p.country_iso
        LIMIT 1
    ), '') AS country_iso2,
    r.position, r.score, r.outcome
FROM fixture_participants fp
JOIN participants p ON p.id = fp.participant_id
LEFT JOIN results r ON r.fixture_id = fp.fixture_id AND r.participant_id = p.id
WHERE fp.fixture_id = ?
ORDER BY fp.role;
