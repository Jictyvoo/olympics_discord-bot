# Errors

**Errors are values.** Functions return `error` as their last result. `panic` is reserved for
unrecoverable boot-time failures inside `cmd/<binary>/main.go`.

## Wrap with `%w`

When forwarding an error across a package boundary, wrap with context:

```go
fixture, err := repo.GetFixture(ctx, id)
if err != nil {
    return fmt.Errorf("syncer: get fixture %s: %w", id, err)
}
```

- Use `%w`, never `%v` or `%s`, so `errors.Is` and `errors.As` keep working.
- Lead with the action you tried.
- The wrapper supplies the trailing `: <cause>` when printed.

## Sentinels

```go
var ErrNotImplemented = errors.New("provider: not implemented")
```

- Exported package vars when callers need to branch on them.
- Match with `errors.Is(err, ErrFoo)`, never with `err == ErrFoo`.
- Don't add fields to a sentinel — promote to a typed error if you need context.

## Typed errors

```go
type MigrationError struct {
    Version string
    Err     error
}
func (e *MigrationError) Error() string { return fmt.Sprintf("apply %s: %v", e.Version, e.Err) }
func (e *MigrationError) Unwrap() error { return e.Err }
```

Match with `errors.As(err, &MigrationError{})`.

## `errors.Join`

Use `errors.Join` when accumulating independent errors across a loop:

```go
var errs []error
for _, p := range s.providers.Enabled() {
    if err := s.syncProviderDay(ctx, p, day); err != nil {
        errs = append(errs, err)
    }
}
return errors.Join(errs...)
```

## Forbidden

- `github.com/pkg/errors` — removed from this repo.
- String comparison on `err.Error()` to detect a condition.
- Swallowing errors with `_ = err`. Either handle, log+return, or document why with a one-line
  comment.
- `log.Fatal*` outside `main` — return the error and let the caller decide.
- `panic` outside `main` after DI failure.

## At the binary boundary

`main()` is the only place that converts an error into an exit code:

```go
if err := serve(); err != nil {
    slog.Error("serve failed", slog.String("err", err.Error()))
    os.Exit(1)
}
```

No `os.Exit` anywhere else.
