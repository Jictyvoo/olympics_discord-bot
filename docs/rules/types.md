# Types

## Domain types live in `internal/domain/eventcore/`

- `Fixture`, `Participant`, `Result`, `Notification`, `DiscordEvent`, … — provider-neutral.
- Provider-specific concepts (medal, FIFA group code, Olympic discipline) do not become fields on
  these structs. They get normalized into the existing fields (`Outcome`, `Stage.Name`, etc.) at
  the fetcher boundary.

## Canonical IDs

- `CanonicalID [16]byte` — deterministic from `(ProviderID, ExternalKey)` via `pkg/idgen.New`.
- Stored as 32-char lowercase hex in the database.
- Never store raw provider IDs; always pass through `idgen.New`.
- `ExternalID{Provider, Key}` is the addressing pair — keep both together.

## Typed strings instead of enums

Go has no enums. Use typed strings with a `Valid()` method:

```go
type FixtureStatus string

const (
    FixtureStatusScheduled FixtureStatus = "scheduled"
    FixtureStatusLive      FixtureStatus = "live"
    FixtureStatusFinished  FixtureStatus = "finished"
    FixtureStatusCancelled FixtureStatus = "cancelled"
)

func (s FixtureStatus) Valid() bool { … }
```

Fetcher mappers translate upstream strings into these constants at the boundary; domain code
trusts the constant set.

## Nullable columns

- `*int`, `*time.Time`, `*string` for genuinely nullable domain fields.
- Empty string is not "null" — distinguish explicitly when needed.
- Helpers in `internal/infra/repositories/repocommon/` convert between SQL nullable types and
  domain pointers; never sprinkle the conversions across repo code.

## Value vs pointer ownership

- `CanonicalID`, typed strings, `Outcome` → value semantics, copies are cheap.
- `Fixture`, `Participant`, `Result` → value-passed when read-only; returned by value from repos.
- `*FixtureRepo`, `*Notifier`, `*Syncer` → pointer; they hold connections / channels / locks.

## Errors

- Sentinel errors are `var ErrFoo = errors.New("…")` at package level.
- Error types implement `Unwrap() error` if they need to nest a cause.
- Don't add fields to sentinel errors — use a typed error (`type FooError struct{…}`) and an `Is`
  method.

## Time

- All `time.Time` values stored or compared in UTC. Convert at the boundary (`t.UTC()`).
- Durations come from config; never `time.Duration` literals in business code.
- `time.Time{}` is the only acceptable "zero" sentinel — never `time.Unix(0,0)`.

## What NOT to do

- No naked maps (`map[string]any`) in domain types. Define a struct.
- No `any` as a domain field.
- No mutable global state. If you find a `var notifier = …`, it's a bug.
- No type aliases to flatten import chains (`type Foo = otherpkg.Foo`) outside the
  generated `dbgen` stub. They obscure where the type really lives.
