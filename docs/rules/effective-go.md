# Effective Go — enforced subset

Project-specific subset of the `effective_go` style guide. Everything here is enforced by review
or by `golangci-lint`.

## Receivers

- Decide value vs pointer once per type and stick with it for every method on that type.
- Default to pointer receivers when the type holds a `sync.Mutex`, holds a connection/file
  handle, is larger than ~64 bytes, or is mutated by any method.
- Small immutable value types use value receivers (`CanonicalID`, `Outcome`).

## Initialization

- No `init()` functions that do I/O, read environment, or open network/database connections.
- Compile-time `init()` registering blank-imported drivers (`_ "modernc.org/sqlite"`) is allowed
  inside `cmd/<binary>/` only.
- Package-level vars: zero-value initialization or simple literal expressions only. Anything that
  could fail belongs in a constructor.

## Error handling

- Lead with the happy path; early-return on errors.
- Don't shadow `err` inside nested scopes — name the inner error or chain explicitly.
- No naked `return`. No `if err == nil` immediately before a return.

```go
// good
fixture, err := repo.GetFixture(ctx, id)
if err != nil {
    return fmt.Errorf("get fixture %s: %w", id, err)
}
return notifier.Notify(ctx, fixture)
```

## Loops and ranges

- `for i, v := range xs` only when both are used.
- `for _, v := range xs` when the index is unused.
- `for range ch` when only iteration matters.
- Prefer `slices.SortFunc`, `slices.Contains`, `maps.Keys` over hand-rolled loops.

## Slices and maps

- Pre-size with `make([]T, 0, n)` when `n` is known.
- Return `nil` slices, not empty slices, for "no results" — unless an API contract requires `[]T{}`.
- Maps as parameters: callee may not mutate; pass a copy if mutation is needed.

## Formatting

- `gofumpt` formatting, enforced by `task tools:fmt`.
- `goimports` import order, enforced by golangci-lint.
- 100-column soft limit; `golines` breaks longer lines via the formatter.

## Comments

- Default: no comment. Identifiers carry the meaning.
- Write a comment only for non-obvious business rules, subtle invariants, or workarounds for a
  specific external constraint.
- Never describe WHAT the code does. Never reference the current PR / task / ticket inline.
- Package doc comment (`// Package foo …`) is required for exported packages used outside the
  immediate parent.

## What NOT to do

- No multi-paragraph docstrings or multi-line `/* … */` blocks. One short line, max.
- No `interface{ … }` declarations inline at a call site — name them.
- No empty branches (`if err != nil { /* ignore */ }`). Either handle or `_ = err` with a
  one-line comment explaining why.
- No `any` outside boundary code (JSON decoding, generic helpers). Domain code is typed.
- No `time.Now()` inside use cases — inject a clock, or take the time as a parameter.
