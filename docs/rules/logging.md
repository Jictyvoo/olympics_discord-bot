# Logging

**`log/slog` only.** No `fmt.Println`, no `log.Printf`, no third-party loggers. No wrapper
package. The `forbidigo` lint rule blocks `fmt.Print*` outside test files.

## Call shape

- The package-level `slog.Info` / `slog.Error` is acceptable at the very top of `main()`.
- Inside packages, take a `*slog.Logger` as a constructor argument or function parameter. Never
  reach for a global.
- Always use typed attributes:

```go
logger.Info("fixture upserted",
    slog.String("fixture_id", fixture.ID.String()),
    slog.String("provider", string(fixture.Ext.Provider)),
)
```

- Errors: `slog.String("err", err.Error())`. Don't pass the error directly — the default
  formatter is noisy.

## Handler setup

`cmd/<binary>/main.go` configures the handler once:

```go
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: levelFromConfig(conf),
}))
slog.SetDefault(logger)
```

- Production: JSON handler.
- Development (`conf.Debug == true`): Text handler with `AddSource: true`.
- Level threshold from config; no per-package logger overrides.

## Levels

- `Debug` — high-volume internal state useful only when investigating.
- `Info` — lifecycle events (boot, sync tick complete, fixture upserted).
- `Warn` — degraded but continuing (one provider failed; another succeeded).
- `Error` — operation failed; the caller must handle.

Never use `Fatal` (slog doesn't have it for a reason).

## Context

When `ctx` carries useful values (request ID, sync cycle ID), attach them on the logger:

```go
logger := logger.With(slog.String("sync_cycle", cycleID))
```

No hidden globals; the enriched logger is passed forward by parameter.

## Forbidden

- `fmt.Println`, `fmt.Printf`, `log.Print*` outside `_test.go` files.
- Any wrapper that hides the slog API.
- Format-string-style logging: `slog.Info("fixture %s failed", id)`.
- Mutating the default logger after `main` finishes setup.
