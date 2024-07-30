package config

import (
	"os"
	"strconv"
)

const (
	envUseDebug        = "DEBUG"
	envProjectName     = "PROJECT_NAME"
	envDiscordToken    = "DISCORD_TOKEN"
	envDiscordClientID = "DISCORD_CLIENT_ID"

	envDatabasePath   = "DATABASE_PATH"
	envWatchCountries = "WATCH_COUNTRIES"
	envAPILocale      = "API_LOCALE"
)

type EnvLoader struct{}

func (l EnvLoader) GetString(element *string, name string) string {
	if result := os.Getenv(name); result != "" {
		*element = result
	}
	return *element
}

func (l EnvLoader) GetUint16(element *uint16, name string) uint16 {
	resultStr := os.Getenv(name)

	if value, err := strconv.ParseUint(resultStr, 10, 16); err == nil {
		*element = uint16(value)
	}
	return *element
}
