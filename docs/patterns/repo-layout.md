# Repository layout

Each SQL dialect has the same structure under `internal/infra/repositories/<dialect>/`:

```
internal/infra/repositories/<dialect>/
├── dbgen/            sqlc-generated; never edited by hand
│   ├── models.go     row structs
│   └── *.sql.go      query funcs
├── queries/
│   ├── fixtures.sql
│   ├── participants.sql
│   └── …
├── repo.go           shared base struct (repoSQLite / repoMySQL) + withTimeout
├── <entity>_repo.go  one repo per entity; wraps dbgen.Queries
└── di.go             Register(inj, db) — concrete + interface bindings
```

Shared helpers live in `internal/infra/repositories/repocommon/`.

## Adding a column

1. Edit `build/entschema/<Entity>.go`.
2. Run `task gen:migration DIFF_NAME=add_<entity>_<col>` — Atlas writes a new `*.sql` migration
   into `build/migrations/<dialect>/`.
3. Add the column to the matching `queries/*.sql` SELECT / INSERT / UPDATE.
4. Run `task gen:code` — sqlc regenerates `dbgen/`.
5. Update the `rowTo<Entity>` mapper in the repo file.
6. Add the field to the matching `eventcore` type if the column maps to a domain concept.

## Mapper pattern

Each `<entity>_repo.go` owns `rowTo<Entity>` and `<entity>ToParams` converters. Use
`repocommon` helpers:

```go
repocommon.NullStr(row.SomeField)    // any → string ("" if nil)
repocommon.NullInt(row.Position)     // any → *int
repocommon.NullTime(row.SentAt)      // any → *time.Time
repocommon.TimeOrZero(row.SentAt)    // any → time.Time (zero if nil)
repocommon.IDFromStr(row.ID)         // 32-char hex → CanonicalID
repocommon.IDToStr(id)               // CanonicalID → 32-char hex
```

## `withTimeout` pattern

Every query derives its context from the caller:

```go
func (r *FixtureRepo) GetFixture(ctx context.Context, id eventcore.CanonicalID) (eventcore.Fixture, error) {
    qctx, cancel := withTimeout(ctx)
    defer cancel()
    row, err := r.queries.GetFixture(qctx, repocommon.IDToStr(id))
    if err != nil {
        return eventcore.Fixture{}, err
    }
    return rowToFixture(row), nil
}
```

Never call `context.Background()` inside a repo — always derive from the caller's `ctx`.

## DI — concrete + interface

In `<dialect>/di.go`, register both the concrete `*Repo` and the consumer interface binding:

```go
remy.RegisterConstructorArgs1(inj, remy.Factory[*FixtureRepo], NewFixtureRepo)
remy.RegisterConstructorArgs1(inj, remy.Factory[syncer.FixtureWriter],
    func(r *FixtureRepo) syncer.FixtureWriter { return r })
```

Use-case packages depend only on the interface; `bootstrap` may resolve the concrete repo when
it needs it (e.g. for seeding).

## Per-dialect SQL differences

| Concern       | SQLite                                             | MySQL                                                  |
|---------------|----------------------------------------------------|--------------------------------------------------------|
| Upsert        | `INSERT … ON CONFLICT (…) DO UPDATE` + `RETURNING` | `INSERT … ON DUPLICATE KEY UPDATE` + separate `SELECT` |
| Now()         | `DATETIME('now')`                                  | `NOW()` or `CURRENT_TIMESTAMP`                         |
| Autoincrement | `INTEGER PRIMARY KEY AUTOINCREMENT`                | `BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY`           |
| Boolean       | `INTEGER` 0/1                                      | `TINYINT(1)`                                           |

The query SQL is duplicated per dialect. Schema source of truth (Ent) is single — Atlas emits
the dialect-specific DDL.

## MySQL stub

While `task gen:code:mysql` hasn't been run against a real DDL, `repomysql/dbgen/` is a
stub package aliasing `reposqlite/dbgen` types and returning `errNotImplemented` from every
method. The repo files in `repomysql/` compile against the stub and become live the moment sqlc
regenerates.
