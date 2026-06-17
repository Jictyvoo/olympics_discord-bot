# Security

## Secrets

- Bot tokens, database DSNs, API keys come from **environment variables only**.
- The corresponding `Config` field is tagged `toml:"-"` so a misconfigured TOML cannot commit
  them by accident.
- Never log a token, a full DSN, or any secret-bearing struct. When you must log something near a
  secret, log a hash prefix or a redacted form.
- `git ls-files | grep -i 'token\|secret\|password'` should return only docs and tests.

## TLS

- HTTPS for every upstream. Plain HTTP allowed only for localhost in dev.
- `InsecureSkipVerify` is forbidden. If a cert chain doesn't validate, fix the chain.
- The default `httpdatasource` client uses `MinVersion: tls.VersionTLS12`.

## Untrusted input

- Treat every HTTP response body as untrusted JSON. Decode into a typed struct, never `any`.
- Validate at the boundary: status codes, expected `Content-Type`, field presence.
- Never log raw response bodies at `Info` — they may contain user PII. Bodies belong at `Debug`.

## SQL

- All queries go through sqlc-generated `dbgen` code → parameter binding is automatic.
- No `fmt.Sprintf` to build a query string.
- No raw `db.Exec(rawSQL)` outside the migrator.

## Discord

- The bot token is held by the `*discordgo.Session` injected into `discordfacade`. No package
  outside `discordfacade` may read or hold a reference to the token.
- Discord callbacks are validated for guild membership before acting.

## Crypto

- `crypto/sha256` for the fixture checksum. Not for anything else without review.
- IDs are random or deterministic-by-design (`pkg/idgen`); never use timestamps as IDs.

## Dependency hygiene

- `task tools:vuln` runs `govulncheck`. Run before every release.
- `go mod tidy` produces no diff in CI.
- Pin tool versions in `tools/go.mod`, separate from runtime deps.

## What NOT to do

- No `init()` that mutates global state — makes audit harder.
- No `unsafe` package usage.
- No silent error swallowing on a security path (auth, signature verify, token refresh).
- No "TODO security" comments. Either fix or open an issue.
