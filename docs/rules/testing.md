# Testing

## Colocation

- Unit tests live next to the source: `fixture.go` → `fixture_test.go` in the same package.
- Mocks live next to the interface they implement: `dispatcher.go` → `dispatcher_mock_test.go`
  in the same package. `package <name>` (test-only because of the `_test.go` suffix).
- Integration tests are guarded by `//go:build integration` and run with
  `go test -tags=integration ./...`.

## Table-driven by default

```go
func TestFixtureChecksum(t *testing.T) {
    testCases := []struct {
        name string
        in   eventcore.Fixture
        want string
    }{
        {"empty", eventcore.Fixture{}, "e3b0c4…"},
        {"with participants", fxWithParts(), "abc123…"},
    }
    for _, tCase := range testCases {
        t.Run(tCase.name, func(t *testing.T) {
            got := tCase.in.Checksum()
            if got != tCase.want {
                t.Fatalf("got %q want %q", got, tCase.want)
            }
        })
    }
}
```

- Always iterate with `tCase` (not `tc`, not `c`).
- Always `t.Run(tCase.name, …)` so failures point at the case.
- No shared mutable state between cases.

## Mocks

- Generated or hand-rolled — both fine. If generated, regenerate via `task tools:mocks`.
- The mock is `package <same as interface>`, so it can satisfy unexported interfaces without
  visibility tricks.
- Name: `<TypeName>Mock`. Methods record calls and return canned values.
- See [`docs/patterns/mocks.md`](../patterns/mocks.md) for the full recipe.

## What unit tests must NOT do

- Touch the filesystem (except `t.TempDir()`).
- Open sockets.
- Spawn subprocesses.
- Read environment variables (except via a test-injected `Config`).
- Use `time.Now()` — inject a clock instead.

## Integration tests

- Build tag: `//go:build integration` on the test file.
- Skip if the required environment isn't present:

```go
//go:build integration

func openTestDB(t *testing.T) *sql.DB {
    dsn := os.Getenv("TEST_MYSQL_DSN")
    if dsn == "" { t.Skip("TEST_MYSQL_DSN not set") }
    …
}
```

- One package per integration concern (`repomysql_integration_test.go` etc.).
- Never depend on a live external API. Replay cached responses (see `.rest_cache/` pattern) or
  spin up a local test double.

## Coverage targets

- Domain packages (`internal/domain/`): high — they're pure logic.
- Infra packages: target the mapping layer (`rowTo<T>`, DTO → eventcore) and the cache
  decorator. Don't unit-test the SQL itself; that's what integration tests are for.
- Bootstrap, main: smoke-tested via `go build` + a single startup test that resolves the top
  service.

## What NOT to do

- No global test fixtures via `TestMain`. Each test sets up what it needs.
- No `time.Sleep` to wait for goroutines. Use synchronisation.
- No assertions hidden behind helper packages that obscure the failure site.
- No `t.Parallel()` unless the test actually demonstrates a parallel concern; default is
  sequential.
