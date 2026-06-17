# olhojogo

Multi-provider sports event monitoring daemon. Pulls schedules and results from upstream sports
APIs (Olympics SPH, FIFA, …), persists them to SQLite or MySQL, and emits Discord notifications
and scheduled events.

## Documentation map

- **[architecture.md](architecture.md)** — high-level data flow and package responsibilities.
- **[rules/](rules/)** — mandatory rules for all new and modified code.
- **[patterns/](patterns/)** — implementation recipes for common tasks.

## Binaries

- `cmd/olhojogo/` — main daemon.
- `cmd/devfetch/` — debug CLI for one-shot fetches and cache replays.

## Build

```bash
task build:bin              # CGO-free; outputs ./bin/olhojogo and ./bin/devfetch
task tools:lint             # golangci-lint
task test                   # go test ./...
podman build -f Containerfile -t olhojogo .
```

## Config

```bash
cp conf.example.toml conf.toml
export OLH_DISCORD_TOKEN=...
./bin/olhojogo serve
```

See [rules/config.md](rules/config.md) for the full schema and env var conventions.
