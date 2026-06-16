package config

// Environment variable names used by LoadFromEnv.
// All share the OLH_ prefix; rename to SPORTSMON_ at the module rename commit.
const (
	envDebug     = "OLH_DEBUG"
	envDBDriver  = "OLH_DB_DRIVER"
	envDBDSN     = "OLH_DB_DSN"
	envDBMigrate = "OLH_DB_MIGRATE"
	envCacheBack = "OLH_CACHE_BACKEND"
	envCachePath = "OLH_CACHE_PATH"
	envCacheTTL  = "OLH_CACHE_TTL"
	envDiscToken = "OLH_DISCORD_TOKEN" //nolint:gosec // env var name, not a credential
	envDiscGuild = "OLH_DISCORD_GUILD_ID"
	envDiscChan  = "OLH_DISCORD_CHANNEL"
	envSyncInt   = "OLH_SYNC_INTERVAL"
	envNotifWin  = "OLH_NOTIFY_WINDOW"
	envDiscHoriz = "OLH_DISCORD_HORIZON"
)
