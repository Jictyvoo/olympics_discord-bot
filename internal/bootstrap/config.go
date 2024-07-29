package bootstrap

import (
	"log/slog"
	"os"

	"github.com/jictyvoo/olympics_data_fetcher/pkg/config"
)

func Config() config.Config {
	conf, confErr := config.LoadTOML(config.DefaultFileName)
	if confErr != nil {
		slog.Error(
			"Error loading config file",
			slog.String("file", config.DefaultFileName),
			slog.String("error", confErr.Error()),
		)
		os.Exit(1)
	}

	config.LoadConfigFromLoader(&conf, config.EnvLoader{})
	return conf
}
