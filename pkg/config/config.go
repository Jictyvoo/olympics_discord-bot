package config

const DefaultFileName = "conf.toml"

type (
	Runtime struct {
		WatchCountries []string `toml:"watch_countries"`
		APILocale      string   `toml:"api_locale"`
	}

	Discord struct {
		Token    string `toml:"token"`
		ClientID string `toml:"client_id"`
	}

	Config struct {
		IsDebug     bool    `toml:"is_debug"`
		ProjectName string  `toml:"project_name"`
		Runtime     Runtime `toml:"server"`
		Discord     Discord `toml:"discord"`
	}
)
