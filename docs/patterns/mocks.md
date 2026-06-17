# Colocated mocks

Mocks live in the same Go package as the interface they implement, in a file named
`<source>_mock_test.go`. The `_test.go` suffix makes them test-only — they don't ship in the
production binary.

## Why colocated

- The mock can implement an unexported interface without visibility tricks.
- Touching the interface and its mock means editing two files in the same directory — easy to
  notice if one drifts.
- No `internal/domain/mocks/` package to maintain or audit.
- Failures point at the right package; no cross-package import in test output.

## File naming

| Source                                                     | Mock                      |
|------------------------------------------------------------|---------------------------|
| `dispatcher.go` (declares `Dispatcher`)                    | `dispatcher_mock_test.go` |
| `interfaces.go` (declares `FixtureWriter`, `ResultWriter`) | `interfaces_mock_test.go` |
| `fetcher.go` (declares `httpClient`)                       | `fetcher_mock_test.go`    |

One mock file per source file. If the source has multiple interfaces, the mock file holds them
all.

## Hand-rolled mock template

```go
// dispatcher_mock_test.go
package notifier

import "context"

type dispatcherMock struct {
    sendCalls []struct {
        channelID string
        content   string
    }
    sendMessageID string
    sendErr       error
}

func (m *dispatcherMock) Send(ctx context.Context, channelID, content string) (string, error) {
    m.sendCalls = append(m.sendCalls, struct {
        channelID string
        content   string
    }{channelID, content})
    return m.sendMessageID, m.sendErr
}
```

Record every call site for verification. Stub returns are bare struct fields, not channels —
tests should be deterministic.

## Generated with mockgen

`tools/Taskfile.yml` exposes a `mocks` target:

```bash
task tools:mocks
```

Under the hood:

```bash
go tool -modfile=tools/go.mod mockgen \
    -source=internal/domain/usecases/notifier/dispatcher.go \
    -destination=internal/domain/usecases/notifier/dispatcher_mock_test.go \
    -package=notifier
```

Add one `mockgen` invocation per source file that declares an interface needing a mock.

## Using mocks in tests

```go
func TestNotifier_NotifyPending(t *testing.T) {
    disp := &dispatcherMock{sendMessageID: "msg-1"}
    n := New(fixtureReaderMock{…}, notificationRepoMock{…}, disp, markdownRenderer, "chan", time.Hour)
    err := n.NotifyPending(context.Background())
    if err != nil { t.Fatal(err) }
    if len(disp.sendCalls) != 1 { t.Fatalf("expected 1 Send, got %d", len(disp.sendCalls)) }
}
```

The mock is package-local, so the test can construct it directly — no setter ceremony, no
interface assertion.

## Forbidden

- A separate `mocks/` package.
- Mocks under a `testing/` build tag (use the `_test.go` suffix instead).
- Mocks that import the production code's third-party deps. Mocks are dumb data + bare returns.
- Mocks that contain logic. If your mock needs an `if`, your test is testing the mock.
