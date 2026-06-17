# HTTP clients

All outbound HTTP goes through `internal/infra/httpdatasource/`. No `http.Get`, no top-level
`http.DefaultClient`, no per-package HTTP setup.

## The `Client` interface

```go
type Request struct {
    Method, URL string
    Headers     map[string]string
    Body        []byte
    Timeout     time.Duration
    CacheKey    string // optional; "" = no cache
}

type Response struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

type Client interface {
    Do(ctx context.Context, req Request) (Response, error)
}
```

The default implementation uses `net/http` with TLS 1.2+ and per-request timeouts. CGO-free.

## Cache integration

Cache decisions are at the call site, not in the client. The fetcher decides whether to
consult `cachestore.Cache` before issuing the request:

```go
if cached, ok, _ := f.cache.Get(ctx, cacheKey); ok {
    return parse(cached)
}
resp, err := f.http.Do(ctx, req)
…
_ = f.cache.Set(ctx, cacheKey, resp.Body, ttl)
```

The `httpdatasource` package itself never caches. Two layers, two responsibilities.

## Fallback to libcurl

The `andelf/go-curl` path is preserved behind `//go:build curl` in
`internal/infra/httpdatasource/curl_client.go`. Build with `go build -tags curl ./cmd/olhojogo`
when an upstream needs libcurl's TLS quirks.

Default builds do not link libcurl and do not require CGO.

## Timeouts and retries

- Every `Do` call respects `ctx` cancellation and a per-request `Timeout`.
- The client does not retry. Retries with backoff are the caller's responsibility (typically a
  fetcher) so the policy stays visible.
- Default request timeout: 30 seconds. Override per-call when an upstream is known to be slow.

## What NOT to do

- No `http.DefaultClient`.
- No package-level `var client = &http.Client{…}` outside `httpdatasource`.
- No HTTP calls outside `infra/fetchers/` or `infra/discordfacade/` (which uses `discordgo`'s
  internal HTTP, not `httpdatasource`).
- No TLS skip-verify. Ever. If a cert chain fails, fix the chain.
