# Config layout

`pkg/confloader` provides a generic, type-safe TOML + env loader. `config/` declares the typed
`Config` struct, env var name constants, and per-domain binders.

## File layout

```
config/
├── config.go      Config + sub-struct types + Default()
├── env_vars.go    const env* = "OLH_..."
└── load.go        Load() + loadFromEnv() + per-domain binders
```

## Config struct

```go
package config

type Config struct {
    ProjectName string        `toml:"project_name"`
    Debug       bool          `toml:"debug"`
    Runtime     Runtime       `toml:"runtime"`
    Database    Database      `toml:"database"`
    Cache       Cache         `toml:"cache"`
    Discord     Discord       `toml:"discord"`
    Providers   []ProviderCfg `toml:"providers"`
}

type Runtime struct {
    SyncInterval   time.Duration `toml:"sync_interval"`
    NotifyWindow   time.Duration `toml:"notify_window"`
    DiscordHorizon time.Duration `toml:"discord_horizon"`
}

type Database struct {
    Driver        string        `toml:"driver"`
    DSN           string        `toml:"-"` // env only
    RunMigrations bool          `toml:"run_migrations"`
    Timeout       time.Duration `toml:"timeout"`
}

type Discord struct {
    Token          string `toml:"-"` // env only
    GuildID        string `toml:"guild_id"`
    DefaultChannel string `toml:"default_channel"`
}
```

Secrets are tagged `toml:"-"` so a misconfigured TOML cannot leak them.

## Default()

```go
func Default() Config {
    return Config{
        ProjectName: "olhojogo",
        Runtime: Runtime{
            SyncInterval:   4 * time.Minute,
            NotifyWindow:   3*time.Hour + 30*time.Minute,
            DiscordHorizon: 14 * 24 * time.Hour,
        },
        Database: Database{Driver: "sqlite", Timeout: 30 * time.Second},
        Cache:    Cache{Backend: "file", FilePath: ".cache", TTL: 4 * time.Minute},
    }
}
```

`Default` returns a zero-value-safe struct. TOML overrides it; env overrides TOML.

## Load

```go
func Load(path string) (Config, error) {
    return confloader.Load(path, Default(), loadFromEnv)
}

func loadFromEnv(c *Config) error {
    return confloader.BindEnv(
        confloader.BindField(&c.Debug,              envDebug,         confloader.ParseBool),
        confloader.BindField(&c.Database.Driver,    envDBDriver,      confloader.ParseString),
        confloader.BindField(&c.Database.DSN,       envDBDSN,         confloader.ParseString),
        confloader.BindField(&c.Database.RunMigrations, envDBMigrate, confloader.ParseBool),
        confloader.BindField(&c.Discord.Token,      envDiscordToken,  confloader.ParseString),
        // …
    )
}
```

## Env var name constants

`config/env_vars.go`:

```go
package config

const (
    envDebug         = "OLH_DEBUG"
    envDBDriver      = "OLH_DB_DRIVER"
    envDBDSN         = "OLH_DB_DSN"
    envDBMigrate     = "OLH_DB_MIGRATE"
    envCacheBackend  = "OLH_CACHE_BACKEND"
    envCachePath     = "OLH_CACHE_PATH"
    envCacheTTL      = "OLH_CACHE_TTL"
    envDiscordToken  = "OLH_DISCORD_TOKEN"
    envDiscordGuild  = "OLH_DISCORD_GUILD_ID"
    envDiscordChan   = "OLH_DISCORD_CHANNEL"
    envSyncInterval  = "OLH_SYNC_INTERVAL"
    envNotifyWindow  = "OLH_NOTIFY_WINDOW"
    envDiscordHoriz  = "OLH_DISCORD_HORIZON"
)
```

Never inline literal env strings in `loadFromEnv`.

## confloader package

`pkg/confloader/` is project-agnostic:

- `Load[T](path, defaults T, envBinder func(*T) error) (T, error)` — TOML + env overlay.
- `BindField[T](*T, env string, parser func(string) (T, error)) Binder` — generic, type-safe.
- `BindEnv(binders ...Binder) error` — accumulates via `errors.Join`.
- Built-in parsers: `ParseBool`, `ParseString`, `ParseInt`, `ParseUint`, `ParseFloat64`, `ParseDuration`.

## Forbidden

- `os.Getenv` outside `config/`.
- `time.Duration` literals inside business code (they belong in `Config`).
- Mutating `Config` after `Load()` returns.
- Reading `Config` from inside a constructor argument that wasn't passed by `bootstrap`.
