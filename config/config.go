package config

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const (
	defaultSyncIntervalMinutes = 4
	defaultDBTimeoutSeconds    = 30
	defaultCacheTTLMinutes     = 4
	defaultDiscordHorizonDays  = 14
	hoursPerDay                = 24
)

type Runtime struct {
	SyncInterval   time.Duration `toml:"sync_interval"`
	NotifyWindow   time.Duration `toml:"notify_window"`
	DiscordHorizon time.Duration `toml:"discord_horizon"`
}

type Database struct {
	Driver        string        `toml:"driver"`
	DSN           string        `toml:"dsn"`
	RunMigrations bool          `toml:"run_migrations"`
	Timeout       time.Duration `toml:"timeout"`
}

type Cache struct {
	Backend  string        `toml:"backend"`
	FilePath string        `toml:"file_path"`
	TTL      time.Duration `toml:"ttl"`
}

type Discord struct {
	Token          string `toml:"token"`
	GuildID        string `toml:"guild_id"`
	DefaultChannel string `toml:"default_channel"`
}

type ProviderCfg struct {
	Code             eventcore.ProviderID `toml:"code"`
	Enabled          bool                 `toml:"enabled"`
	BaseURL          string               `toml:"base_url"`
	Language         string               `toml:"language"`
	WatchedCountries []string             `toml:"watched_countries"`
	DiscordChannel   string               `toml:"discord_channel"`
}

type Config struct {
	ProjectName string        `toml:"project_name"`
	Debug       bool          `toml:"debug"`
	Runtime     Runtime       `toml:"runtime"`
	Database    Database      `toml:"database"`
	Cache       Cache         `toml:"cache"`
	Discord     Discord       `toml:"discord"`
	Providers   []ProviderCfg `toml:"providers"`
}

func Default() Config {
	return Config{
		ProjectName: "olhojogo",
		Runtime: Runtime{
			SyncInterval:   defaultSyncIntervalMinutes * time.Minute,
			NotifyWindow:   3*time.Hour + 30*time.Minute,
			DiscordHorizon: defaultDiscordHorizonDays * hoursPerDay * time.Hour,
		},
		Database: Database{
			Driver:  "sqlite",
			DSN:     "olhojogo.db",
			Timeout: defaultDBTimeoutSeconds * time.Second,
		},
		Cache: Cache{
			Backend:  "file",
			FilePath: ".cache",
			TTL:      defaultCacheTTLMinutes * time.Minute,
		},
		Discord: Discord{
			DefaultChannel: "sports",
		},
		Providers: []ProviderCfg{
			{
				Code:     eventcore.ProviderOlympics,
				Enabled:  true,
				Language: "ENG",
			},
		},
	}
}
