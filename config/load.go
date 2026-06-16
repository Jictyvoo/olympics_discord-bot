package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jictyvoo/olhojogo/pkg/confloader"
)

// Load reads the TOML file at path (falling back to CONF_FILE env, then
// "conf.toml") and overlays env vars on top.
func Load(path string) (Config, error) {
	conf := Default()
	if path == "" {
		path = os.Getenv("CONF_FILE")
	}
	if path == "" {
		path = "conf.toml"
	}
	if _, err := toml.DecodeFile(path, &conf); err != nil && !os.IsNotExist(err) {
		return conf, err
	}
	if err := LoadFromEnv(&conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func LoadFromEnv(conf *Config) error {
	return confloader.BindEnv(
		confloader.BindField(&conf.Debug, envDebug, confloader.ParseBool),
		confloader.BindField(&conf.Database.Driver, envDBDriver, confloader.ParseString),
		confloader.BindField(&conf.Database.DSN, envDBDSN, confloader.ParseString),
		confloader.BindField(&conf.Database.RunMigrations, envDBMigrate, confloader.ParseBool),
		confloader.BindField(&conf.Cache.Backend, envCacheBack, confloader.ParseString),
		confloader.BindField(&conf.Cache.FilePath, envCachePath, confloader.ParseString),
		confloader.BindField(&conf.Cache.TTL, envCacheTTL, confloader.ParseDuration),
		confloader.BindField(&conf.Discord.Token, envDiscToken, confloader.ParseString),
		confloader.BindField(&conf.Discord.GuildID, envDiscGuild, confloader.ParseString),
		confloader.BindField(&conf.Discord.DefaultChannel, envDiscChan, confloader.ParseString),
		confloader.BindField(&conf.Runtime.SyncInterval, envSyncInt, confloader.ParseDuration),
		confloader.BindField(&conf.Runtime.NotifyWindow, envNotifWin, confloader.ParseDuration),
		confloader.BindField(&conf.Runtime.DiscordHorizon, envDiscHoriz, confloader.ParseDuration),
	)
}
