# Adding a new HTTP fetcher

A "provider" in domain terms is a `provider.Strategy`. The concrete HTTP implementation lives
under `internal/infra/fetchers/<name>fetch/`.

## 1. Create the package

```
internal/infra/fetchers/<name>fetch/
├── fetcher.go        Fetcher struct + Code()/DisplayName()/SyncFixturesByDate()/SyncFixtureResults()
├── client.go         HTTP calls via httpdatasource.Client, cache-aware
├── url.go            endpoint builders (URLs come from config)
├── dto.go            package-local API DTOs; never exported
├── mapper.go         dto → eventcore.SyncDelta
├── di.go             Register(inj, baseURL, lang)
└── doc.go            package comment + endpoint list
```

## 2. Implement provider.Strategy

```go
package <name>fetch

type Fetcher struct {
    http    httpdatasource.Client
    cache   cachestore.Cache
    baseURL string
    lang    string
}

func New(http httpdatasource.Client, cache cachestore.Cache, baseURL, lang string) *Fetcher {
    return &Fetcher{http: http, cache: cache, baseURL: baseURL, lang: lang}
}

func (f *Fetcher) Code() eventcore.ProviderID  { return "<name>" }
func (f *Fetcher) DisplayName() string         { return "<Human Name>" }

func (f *Fetcher) SyncFixturesByDate(ctx context.Context, day time.Time) (eventcore.SyncDelta, error) {
    raw, err := f.fetch(ctx, scheduleURL(f.baseURL, day, f.lang))
    if err != nil {
        return eventcore.SyncDelta{}, err
    }
    var dto apiScheduleResponse
    if err := json.Unmarshal(raw, &dto); err != nil {
        return eventcore.SyncDelta{}, fmt.Errorf("<name>fetch: decode schedule: %w", err)
    }
    return mapSchedule(dto), nil
}

func (f *Fetcher) SyncFixtureResults(ctx context.Context, fx eventcore.Fixture) (eventcore.SyncDelta, error) {
    return eventcore.SyncDelta{}, eventcore.ErrNotImplemented
}
```

Return `eventcore.ErrNotImplemented` for methods you don't ship yet; the syncer loop logs and
moves on.

## 3. Mapper rules

- All mapping from DTO → `eventcore.*` happens in `mapper.go`, never in `fetcher.go`.
- IDs always go through `idgen.New(<name>fetch.ProviderCode, externalKey)`.
- `Outcome` must be one of the constants in `eventcore` (`OutcomeWin`, `OutcomeMedalGold`, …).
- DTOs stay unexported.

## 4. DI registration

```go
// di.go
func Register(inj remy.Injector, baseURL, lang string) {
    remy.RegisterConstructorArgs2(inj, remy.Factory[*Fetcher],
        func(http httpdatasource.Client, cache cachestore.Cache) *Fetcher {
            return New(http, cache, baseURL, lang)
        })
    remy.RegisterConstructorArgs1(inj, remy.Factory[provider.Strategy],
        func(f *Fetcher) provider.Strategy { return f })
}
```

The `provider.Strategy` interface binding is what `provider.Set` consumes.

## 5. Wire in bootstrap

`internal/bootstrap/injections.go`:

```go
case eventcore.Provider<Name>:
    <name>fetch.Register(inj, pc.BaseURL, pc.Language)
```

## 6. Add config

`conf.example.toml`:

```toml
[[providers]]
code     = "<name>"
enabled  = false
base_url = "https://api.example.com"
language = "en"
```

## 7. Checklist before merge

- [ ] `go build ./...` passes with no CGO.
- [ ] `go test ./internal/infra/fetchers/<name>fetch/...` passes.
- [ ] `mapper.go` has table-driven tests for the happy path and one edge case.
- [ ] DTOs are unexported.
- [ ] `doc.go` lists every upstream URL the fetcher hits.
- [ ] No hardcoded URL or language anywhere outside `url.go` / config.
