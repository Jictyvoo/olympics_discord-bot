-- name: UpsertSyncState :exec
INSERT INTO sync_states (provider_id, scope, cursor, last_synced_at, last_error)
VALUES (?, ?, ?, DATETIME('now'), NULL)
ON CONFLICT(provider_id, scope) DO UPDATE SET
    cursor         = excluded.cursor,
    last_synced_at = DATETIME('now'),
    last_error     = NULL;

-- name: GetSyncState :one
SELECT * FROM sync_states WHERE provider_id = ? AND scope = ? LIMIT 1;

-- name: SetSyncStateError :exec
UPDATE sync_states SET last_error = ?, last_synced_at = DATETIME('now')
WHERE provider_id = ? AND scope = ?;
