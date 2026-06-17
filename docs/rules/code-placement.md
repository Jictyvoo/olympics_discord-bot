# Code placement

A new file lives in exactly one of these places. If you're unsure, the layout below answers it.

## Repository tree

```
/
├── Containerfile, Taskfile.yml, go.mod, sqlc.yaml, conf.example.toml
├── cmd/<binary>/                  binary entrypoints; one subdir per executable
├── config/                        typed Config + env binders for both binaries
├── internal/
│   ├── bootstrap/                 DI composition root; picks dialect + enabled fetchers
│   ├── domain/                    business rules; provider-neutral
│   │   ├── eventcore/             canonical domain types (Fixture, Participant, …)
│   │   ├── provider/              Strategy port + Set + SyncDelta
│   │   ├── services/              long-lived behaviour (notifier loop, observers)
│   │   └── usecases/              one package per use case (syncer, notifier, discordsync)
│   ├── infra/                     adapters that satisfy domain ports
│   │   ├── fetchers/              read-only HTTP fetchers (olympicsfetch, fifafetch)
│   │   ├── repositories/          DB-mutating repos (reposqlite, repomysql, repocommon)
│   │   ├── httpdatasource/        outbound HTTP client
│   │   ├── cachestore/            file + memory caches
│   │   └── discordfacade/         bwmarrin/discordgo wrapper
│   └── migrator/                  embedded SQL migration runner
├── pkg/                           small, reusable, project-agnostic libs
│   ├── confloader/                TOML + env loader
│   └── idgen/                     CanonicalID generation
├── build/
│   ├── entschema/                 ent schemas (own go.mod)
│   ├── migrations/{sqlite,mysql}/ embedded SQL migrations
│   └── seed/                      data sources for tools/seedgen
├── tools/                         tools-only module; isolated deps
│   ├── go.mod, tools.go
│   └── seedgen/                   countries.json → migration SQL generator
└── docs/                          rules + patterns + architecture
```

## What lives where

- **`cmd/<binary>/`** — entrypoint and only the entrypoint. Reads config, builds the injector,
  dispatches subcommands, exits. No business logic.
- **`config/`** — typed Config struct + env binders. The only package that calls `os.Getenv`.
- **`internal/bootstrap/`** — the only place that knows about concrete dialect or fetcher choices.
  `switch conf.Database.Driver { … }` lives here, nowhere else.
- **`internal/domain/`** — pure business rules, provider-neutral types, ports declared as
  interfaces. May import other `domain/` packages and `pkg/`. May NOT import `internal/infra/`.
- **`internal/infra/`** — adapters. May import `internal/domain/` to implement its ports. May NOT
  import `internal/bootstrap/` or another binary's `cmd/`.
- **`internal/infra/repositories/`** — packages that mutate DB state (`reposqlite`, `repomysql`).
- **`internal/infra/fetchers/`** — packages that pull data from upstream HTTP APIs (`olympicsfetch`,
  `fifafetch`). They are read-only from the upstream perspective; they don't share a folder with
  DB repositories so the terminology stays clear.
- **`pkg/`** — small libraries with no knowledge of this project's domain. Anyone could import
  them. Right now: `confloader`, `idgen`.
- **`build/`** — code-gen inputs and outputs that ship as build artifacts (entschema, migrations,
  seed data).
- **`tools/`** — anything that runs at build time, never at runtime. Has its own `go.mod`.
- **`docs/rules/`** + **`docs/patterns/`** — rules are mandatory; patterns are copyable recipes.

## Forbidden cross-imports

- `internal/domain/` importing anything under `internal/infra/`.
- Any non-`bootstrap` package switching on `conf.Database.Driver` or `conf.Providers[].Code`.
- Anything outside `config/` calling `os.Getenv`.
- `pkg/` importing anything under `internal/`.
- `internal/infra/repositories/` and `internal/infra/fetchers/` importing each other.

## When in doubt

If a new file doesn't obviously fit one of the dirs above, the design is probably wrong. Stop
and check `docs/patterns/` for an analogous case before inventing a new top-level package.
