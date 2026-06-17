# Config

Single typed `Config` struct, loaded once at startup from a TOML file plus environment variable
overrides via `pkg/confloader`. No `os.Getenv` calls outside `config/`.

## Where config lives

```
config/
├── config.go      # Config + sub-struct types + Default()
├── env_vars.go    # unexported env var name constants (OLH_*)
└── load.go        # Load(), per-domain env binders
```

## Config struct shape

```go
type Config struct {
    ProjectName string        `toml:"project_name"`
    Debug       bool          `toml:"debug"`
    Runtime     Runtime       `toml:"runtime"`
    Database    Database      `toml:"database"`
    Cache       Cache         `toml:"cache"`
    Discord     Discord       `toml:"discord"`
    Providers   []ProviderCfg `toml:"providers"`
}
```

## Loading

```go
func Load(path string) (Config, error) {
    return confloader.Load(path, Default(), loadFromEnv)
}
```

`confloader.Load`:

1. Honours `CONF_FILE` env var if set; otherwise uses the supplied `path` (default `conf.toml`).
2. Treats a missing file as not-an-error (defaults remain).
3. Calls `loadFromEnv` so env vars overlay file values.

## Binding env vars

```go
func loadFromEnv(c *Config) error {
    return confloader.BindEnv(
        confloader.BindField(&c.Debug, envDebug, confloader.ParseBool),
        confloader.BindField(&c.Database.Driver, envDBDriver, confloader.ParseString),
        confloader.BindField(&c.Database.DSN,    envDBDSN,    confloader.ParseString),
        // …
    )
}
```

Env var name constants are in `env_vars.go` — never inline literal strings:

```go
const (
    envDebug        = "OLH_DEBUG"
    envDBDriver     = "OLH_DB_DRIVER"
    envDBDSN        = "OLH_DB_DSN"
    envDiscordToken = "OLH_DISCORD_TOKEN"
)
```

## Secrets

- `Discord.Token` and the database DSN must come from env only.
- Tag with `toml:"-"` so a misconfigured TOML cannot leak credentials.

## Forbidden

- `os.Getenv("...")` outside `config/`.
- Reading config inside business logic — all use cases receive an already-parsed value.
- Mutating `Config` after `Load()` returns.
- Global `var cfg = config.Load()` at package level.
- Switching on `Config.Database.Driver` outside `internal/bootstrap/`.
