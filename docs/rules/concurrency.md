# Concurrency

## Goroutine ownership

Every goroutine has exactly one owner. The owner is responsible for:

- Knowing when the goroutine should stop.
- Providing the cancellation signal (`ctx` or `chan struct{}`).
- Waiting for the goroutine to finish before returning.

If you can't name the owner, don't spawn the goroutine.

## Context first

- Every long-running function takes `ctx context.Context` as its first parameter.
- The deepest call honours `ctx.Done()` and returns `ctx.Err()` when cancelled.
- Don't derive `context.Background()` inside a use-case or repo; chain from the caller's `ctx`.

```go
// good
func (r *FixtureRepo) GetFixture(ctx context.Context, id ID) (Fixture, error) {
    qctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    …
}
```

## Tickers and loops

The sync loop pattern:

```go
ticker := time.NewTicker(interval)
defer ticker.Stop()
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-ticker.C:
        r.tick(ctx)
    }
}
```

- Always `defer ticker.Stop()`.
- Always check `ctx.Done()` in the `select`.
- `r.tick(ctx)` handles errors internally (log+continue) — a single failed tick doesn't kill the
  loop.

## Channels

- Buffered channels only when you have a specific reason for the buffer size. Document it.
- Closer is the sender. Receivers never close.
- Signalling channels use `chan struct{}`, not `chan bool`.

## Locks

- `sync.Mutex` for short critical sections. `sync.RWMutex` only when the read/write ratio is
  measurably skewed.
- Always `defer mu.Unlock()` immediately after `mu.Lock()`. No conditional unlocks.
- Never hold a lock across a channel send/recv or a network call.

## Weak observers

The notifier registry uses `weak.Pointer[Observer]` so observers that go out of scope are
auto-evicted on the next sweep. See [`docs/patterns/observer-weakptr.md`](../patterns/observer-weakptr.md).

Long-lived observers (registered by `bootstrap`) must hold a strong reference at the call site to
stay alive.

## What NOT to do

- No fire-and-forget `go fn()` — every goroutine has an owner and a stop signal.
- No `time.Sleep` in production paths. Use `ctx`-aware waits.
- No naked `select {}` to block forever — block on `ctx.Done()`.
- No sharing of mutable structs across goroutines without a lock or a channel.
