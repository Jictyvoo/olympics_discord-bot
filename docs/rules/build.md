# Build

## Default build

- `CGO_ENABLED=0` — no libcurl, no native SQLite. The default SQLite driver is `modernc.org/sqlite`.
- Go ≥ 1.24 (we use `weak.Pointer`).
- `task build:bin` produces both binaries under `./bin/`.

## Taskfile entry points

- `task` — runs vet + tests (default).
- `task tools:install` — downloads the tools module.
- `task tools:fmt` — gofumpt + goimports + golines via golangci-lint formatters.
- `task tools:lint` — golangci-lint run.
- `task tools:vuln` — govulncheck.
- `task tools:mocks` — regenerate colocated mocks.
- `task build:bin` — produce binaries.
- `task build:image` — produce the container image.
- `task gen:code` — sqlc generate (both engines).
- `task gen:migration` — Atlas migrate diff for SQLite (`DIFF_NAME=foo task build:migration:gen`).
- `task gen:seed` — run `tools/seedgen` against `build/seed/countries.json`.

Per-directory Taskfiles in `tools/` and `build/` are included by the root Taskfile via Task v3
`includes:`.

## Container

```
podman build -f Containerfile -t olhojogo .
```

The Containerfile is at the repository root. There is no `Dockerfile`. The image is multi-stage,
distroless-based, CGO-free.

## Build tags

- `//go:build curl` — opt into the `andelf/go-curl` HTTP fallback. Default builds skip it.
- `//go:build tools` — compile the tools-only file `tools/tools.go`.
- `//go:build integration` — integration tests (skipped by the default `go test`).

## Code generation

- **sqlc** — generates `internal/infra/repositories/{reposqlite,repomysql}/dbgen/`. Not committed
  to git? Up to the repo; current convention commits it for diff visibility.
- **Atlas** — generates SQL files into `build/migrations/{sqlite,mysql}/`. Atlas runs at build
  time only; the runtime applier is `internal/migrator/`.
- **seedgen** — `tools/seedgen` reads `build/seed/countries.json` and writes idempotent SQL
  migrations into both dialect directories.
- **mockgen** — colocated `<file>_mock_test.go` files in the interface's package.

## Forbidden

- `go run` to launch the daemon — use the built binary.
- `task` targets that hide implicit network calls.
- Build-time downloads of binaries (no `curl | sh`). Tools come from `tools/go.mod`.
- CGO unless explicitly required by a build tag.
