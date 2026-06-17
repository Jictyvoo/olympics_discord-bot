# Naming

## Packages

- Lowercase, single word, no underscores, no plurals unless the package is genuinely a collection
  (`fetchers`, `repositories`).
- Neologism-friendly: `confloader`, `idgen`, `httpdatasource`, `cachestore`, `discordfacade`.
- Concrete repo packages take a `repo*` prefix: `reposqlite`, `repomysql`, `repocommon`.
- Concrete HTTP fetcher packages take a `*fetch` suffix: `olympicsfetch`, `fifafetch`.
- `infra` not `infrastructure`. `usecases` not `use_cases`. `discordsync` not `discord-sync`.

## Files

- `snake_case.go` (e.g. `fixture_repo.go`, `event_notifier.go`).
- Test files: `<source>_test.go` colocated with `<source>.go`.
- Mock files: `<source>_mock_test.go` in the same package as the interface (test-only).
- One concern per file. Split when a file passes ~300 lines.

## Types

- Exported types are PascalCase: `FixtureRepo`, `Notifier`, `Strategy`.
- Interfaces named for the role, not the implementer: `FixtureWriter`, not `FixtureRepoInterface`.
- Single-method interfaces take the method name + `er`: `Dispatcher`, `Renderer`, `Reader`.
- Multi-method interfaces describe the role: `Facade`, `Strategy`, `Set`.

## Functions

- Constructors are `New<Type>`: `NewFixtureRepo`, `NewSyncer`. With error: `NewFoo() (*Foo, error)`.
- Boolean queries start with `Is`, `Has`, `Can`: `IsExpired`, `HasParticipants`, `CanNotify`.
- Side-effect functions use verbs: `UpsertFixture`, `SyncFixturesByDate`.
- Internal helpers stay unexported and concise: `rowToFixture`, not `convertDatabaseRowToFixture`.

## Identifiers

- Acronyms keep case: `ID`, `HTTP`, `URL`, `ISO`, `IOC`. So: `userID`, `httpClient`, `iso2`, `iocCode`.
- The receiver name is the type's first letter, lowercased: `func (r *FixtureRepo) …`.
- Context is always `ctx`; errors are always `err`.

## Constants

- `UpperCamelCase` for exported, `lowerCamelCase` for unexported.
- Group related constants in a single `const ( … )` block.
- Env var names: unexported constants in `config/env_vars.go`, format `OLH_<DOMAIN>_<NAME>`.

## What NOT to do

- No Hungarian notation (`pFixture`, `iCount`).
- No `Interface`, `Impl`, `Manager`, `Helper`, `Util` suffixes in type names.
- No `_t` or `_type` suffixes.
- No abbreviations that aren't already common Go idiom (`req`, `resp`, `ctx`, `cfg` are fine;
  `cust`, `prdct`, `mgr` are not).
