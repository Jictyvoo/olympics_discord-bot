# Bootstrap & DI composition

`internal/bootstrap/injections.go` is the **single composition root**. It is the only place
allowed to switch on `conf.Database.Driver` or `conf.Providers[].Code`. Everything else asks for
dependencies through constructor parameters.

## The composition root

```go
func DoInjections(inj remy.Injector, conf appconfig.Config, db *sql.DB) {
    remy.RegisterInstance(inj, conf)
    remy.RegisterInstance(inj, db)

    httpdatasource.Register(inj)
    cachestore.Register(inj, conf.Cache.Backend, conf.Cache.FilePath, conf.Cache.TTL)

    switch conf.Database.Driver {
    case "mysql":
        repomysql.Register(inj, db)
    default:
        reposqlite.Register(inj, db)
    }

    enabled := make([]eventcore.ProviderID, 0, len(conf.Providers))
    for _, pc := range conf.Providers {
        if !pc.Enabled {
            continue
        }
        switch pc.Code {
        case eventcore.ProviderOlympics:
            olympicsfetch.Register(inj, pc.BaseURL, pc.Language)
        case eventcore.ProviderFIFA:
            fifafetch.Register(inj)
        }
        enabled = append(enabled, pc.Code)
    }
    provider.RegisterSet(inj, enabled)

    discordfacade.Register(inj)
    syncer.Register(inj, conf.Runtime.SyncInterval)
    notifier.Register(inj, conf.Discord.DefaultChannel, conf.Runtime.NotifyWindow)
    discordsync.Register(inj, conf.Discord.GuildID, conf.Runtime.DiscordHorizon)

    wireObservers(inj)
}
```

## Per-package `Register` shape

Every package owns a `di.go` with a `Register(inj remy.Injector, …)` function. The function
registers types via `RegisterConstructor*`; it does not branch on config.

```go
// internal/infra/cachestore/di.go
func Register(inj remy.Injector, backend, rootPath string, ttl time.Duration) {
    if backend == "memory" {
        remy.RegisterConstructor(inj, remy.Factory[Cache],
            func() Cache { return memcache.New(ttl) })
        return
    }
    remy.RegisterConstructorErr(inj, remy.Factory[Cache],
        func() (Cache, error) { return filecache.New(rootPath, ttl) })
}
```

A small switch inside a package's `Register` is OK when it picks between two backends in the
same package family (file vs memory cache). A switch that picks between two *different*
packages belongs in `bootstrap`.

## Observer wiring

Observers register at boot, not inside a constructor. `bootstrap.wireObservers` resolves the
notifier and discord-sync and connects them:

```go
func wireObservers(inj remy.Injector) {
    n, err := remy.Get[*notifier.Notifier](inj)
    if err != nil {
        slog.Error("bootstrap: get notifier", slog.String("err", err.Error()))
        return
    }
    ds, err := remy.Get[*discordsync.DiscordSync](inj)
    if err != nil {
        slog.Error("bootstrap: get discordsync", slog.String("err", err.Error()))
        return
    }
    var obs notifier.Observer = ds
    n.RegisterObserver(&obs)

    // Keep a strong ref so the weak.Pointer doesn't evict.
    remy.RegisterInstance(inj, &obs, "discordsync:observer")
}
```

## Resolving the top-level service

`cmd/<binary>/serve.go` resolves the top-level type and runs it:

```go
runner, err := remy.Get[*syncer.Runner](inj)
if err != nil {
    return fmt.Errorf("serve: get runner: %w", err)
}
return runner.Run(ctx)
```

## Forbidden inside `Register` functions

- Branching on `Config` fields outside the `Register`'s own scope.
- `remy.Get[*X](ret)` inside a `RegisterConstructor` callback (the args slots exist for that).
- `RegisterSingleton(func(ret){…})` — opaque to the graph.
- Side effects (file I/O, network) — `Register` builds the graph; constructors do real work.
