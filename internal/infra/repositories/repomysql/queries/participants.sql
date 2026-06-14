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

-- name: ListParticipantsByFixture :many
SELECT p.id, p.created_at, p.updated_at, p.provider_id, p.external_key, p.kind, p.name, p.code, p.country_iso, p.gender
FROM participants p
JOIN fixture_participants fp ON fp.participant_id = p.id
WHERE fp.fixture_id = ?;
