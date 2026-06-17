# Startup

`main` is the only entrypoint that may have side effects. Nothing else may do I/O at import.

## What `main` does, in order

```go
func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
    slog.SetDefault(logger)

    cmd := "serve"
    if len(os.Args) > 1 {
        cmd = os.Args[1]
    }
    switch cmd {
    case "serve":   exitOn(serve())
    case "migrate": exitOn(migrate())
    case "version": printVersion()
    default:        slog.Error("unknown subcommand"); os.Exit(2)
    }
}
```

`serve()` then:

1. Load `Config` from `config.Load`.
2. Open `*sql.DB`.
3. Run migrations if `conf.Database.RunMigrations`.
4. Build the remy injector with `DuckTypeElements: true`.
5. Call `bootstrap.DoInjections(inj, conf, db)`.
6. Optionally open the `discordgo.Session` and register it as an instance.
7. Resolve the top-level service (`*syncer.Runner`) and start it under a signal-cancelled context.

## What `init()` may NOT do

- Open files or sockets.
- Read environment variables.
- Read or decode embedded `//go:embed` data into runtime state. The data may be embedded; binding
  it to a `Config` or runtime struct happens lazily inside a constructor.
- Register global observers or callbacks.

## Blank-import drivers

Allowed only inside `cmd/<binary>/`:

```go
import _ "modernc.org/sqlite"
import _ "github.com/go-sql-driver/mysql"
```

Driver registration via `init()` is the one acceptable side effect, and only in the binary
package.

## Shutdown

- The signal context cancels on `SIGINT` or `SIGTERM`.
- `defer db.Close()` and `defer session.Close()` run after `serve()` returns.
- Long-running goroutines watch `ctx.Done()` and exit cleanly. See [`concurrency.md`](concurrency.md).

## Forbidden

- `os.Exit` outside `main`.
- `panic` outside `main`, unless it's a programmer error that can never occur in a correct build
  (a missing `default` in an enum-style switch).
- Side-effecting top-level `var x = someFn()` declarations in packages outside `cmd/`.
