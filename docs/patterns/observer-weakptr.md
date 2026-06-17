# Observer with weak.Pointer

## Why

A strong-reference observer slice keeps every registered observer alive forever; consumers must
explicitly `Unregister` or the registry leaks. `weak.Pointer[Observer]` lets the GC drop the
observer when its owner releases it; the registry sweeps stale entries on the next `Notify`.

Trade-off: GC decides timing. **Long-lived observers must hold a strong reference at the call
site** to stay alive. `bootstrap.wireObservers` does this by registering the observer as an
instance in the injector.

## Pattern

```go
import "weak"

type Observer interface {
    OnFixture(ctx context.Context, f eventcore.Fixture)
}

type observerSet struct {
    mu    sync.Mutex
    items []weak.Pointer[Observer]
}

func (s *observerSet) Register(o *Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, weak.Make(o))
}

func (s *observerSet) Notify(ctx context.Context, f eventcore.Fixture) {
    s.mu.Lock()
    defer s.mu.Unlock()
    live := s.items[:0]
    for _, w := range s.items {
        p := w.Value()
        if p == nil {
            continue // GC'd — drop silently
        }
        (*p).OnFixture(ctx, f)
        live = append(live, w)
    }
    s.items = live
}
```

## Keeping an observer alive

```go
// internal/bootstrap/injections.go
var obs notifier.Observer = discordSync
notifier.RegisterObserver(&obs)

// Keep a strong reference in the injector so the weak pointer doesn't evict.
remy.RegisterInstance(inj, &obs, "discordsync:observer")
```

`cmd/<binary>/serve.go` doesn't need to do anything extra — the instance lives in the injector
for the duration of the process.

## Requirements

- Go 1.24+ (the `weak` package is stable from 1.24).
- The interface pointer passed to `Register` must be addressable.
- Don't register the same observer twice; the sweep doesn't deduplicate.
