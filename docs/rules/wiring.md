# Wiring — remy DI

We use `github.com/wrapped-owls/goremy-di/remy` with `Config{DuckTypeElements: true}`. The
injector is built once in `cmd/<binary>/serve.go`, populated by `internal/bootstrap/`, and
consulted to resolve the top-level service.

## Rule 1 — `RegisterConstructor*` only

```go
// good
remy.RegisterConstructor(inj, remy.Factory[*Foo], NewFoo)
remy.RegisterConstructorArgs1(inj, remy.Factory[*Bar], NewBar)
remy.RegisterConstructorArgs2Err(inj, remy.Factory[*Baz], NewBaz)
```

Forbidden shape:

```go
// bad — opaque factory; defeats the DI graph
remy.RegisterSingleton(inj, func(ret remy.DependencyRetriever) (*Foo, error) {
    bar, _ := remy.Get[*Bar](ret)
    return NewFoo(bar), nil
})
```

If a constructor has N dependencies, use `RegisterConstructorArgsN[Err]` so remy resolves them
from the graph. Don't call `remy.Get[*X](ret)` inside a factory — it defeats duck-typing and
hides the dependency.

## Rule 2 — consumer-defined interfaces

Each use-case package declares the smallest interface it needs:

```go
// internal/domain/usecases/syncer/interfaces.go
type FixtureWriter interface {
    UpsertFixture(ctx context.Context, f eventcore.Fixture) (eventcore.Fixture, error)
}
```

The use-case `di.go` is one-liner-per-type:

```go
// internal/domain/usecases/syncer/di.go
func Register(inj remy.Injector) {
    remy.RegisterConstructorArgs5(inj, remy.Factory[*Syncer], New)
}
```

remy resolves `FixtureWriter`, `ParticipantWriter`, etc. by duck-typing against whichever repo
package is registered.

## Rule 3 — repo `di.go` registers both concrete and interface bindings

```go
// internal/infra/repositories/reposqlite/di.go
remy.RegisterConstructor(inj, remy.Factory[*dbgen.Queries],
    func() *dbgen.Queries { return dbgen.New(db) })

remy.RegisterConstructorArgs1(inj, remy.Factory[*FixtureRepo], NewFixtureRepo)
remy.RegisterConstructorArgs1(inj, remy.Factory[syncer.FixtureWriter],
    func(r *FixtureRepo) syncer.FixtureWriter { return r })
```

The interface binding is what use cases actually consume. The concrete `*FixtureRepo` binding is
used only by `bootstrap` or by other repos in the same package.

## Rule 4 — bootstrap is the only switchboard

`internal/bootstrap/injections.go` is the only place that chooses between dialects or fetcher
implementations:

```go
switch conf.Database.Driver {
case "mysql":
    repomysql.Register(inj, db)
default:
    reposqlite.Register(inj, db)
}

for _, pc := range conf.Providers {
    if !pc.Enabled { continue }
    switch pc.Code {
    case eventcore.ProviderOlympics:
        olympicsfetch.Register(inj, pc.BaseURL, pc.Language)
    case eventcore.ProviderFIFA:
        fifafetch.Register(inj)
    }
}
```

No `switch driver` inside any `repo*/di.go`. No `switch code` inside any `*fetch/di.go`.

## Rule 5 — `GetWithPairs` over `GetWith(callback)`

When you need to inject something at resolve time:

```go
// good — typed pairs
session, err := remy.GetWithPairs[*discordgo.Session](inj,
    remy.BindInstance(cancelChan),
)

// bad — callback hides what's being added
session, err := remy.GetWith[*discordgo.Session](inj, func(child remy.Injector) error {
    remy.RegisterInstance(child, cancelChan)
    return nil
})
```

## Forbidden

- `RegisterSingleton(func(ret remy.DependencyRetriever)…)` — anywhere.
- `remy.Get[*concreteType](ret)` inside a registered factory.
- Service-locator usage outside `bootstrap` — application code receives dependencies through
  constructor parameters, not by querying the injector.
- Global `var inj = remy.NewInjector(...)`.
