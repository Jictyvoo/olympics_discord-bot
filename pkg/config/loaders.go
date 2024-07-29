package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Loader interface {
	GetString(element *string, name string) string
	GetUint16(element *uint16, name string) uint16
}

func DefaultConfig() Config {
	return Config{
		IsDebug:     false,
		ProjectName: "olympics_PARIS--2024",
		Server: Server{
			Host: "0.0.0.0",
			Port: 8008,
		},
	}
}

func LoadConfigFromLoader(config *Config, loader Loader) {
	// Load debug info
	var useDebug string
	loader.GetString(&useDebug, envUseDebug)
	if !config.IsDebug && useDebug != "" {
		config.IsDebug = strings.ToLower(useDebug) != "false"
	}

	LoadDiscord(&config.Discord, loader)
	loader.GetString(&config.ProjectName, envProjectName)
	loader.GetUint16(&config.Server.Port, envAPIPort)
}

func LoadDiscord(conf *Discord, loader Loader) {
	loader.GetString(&conf.Token, envDiscordToken)
	loader.GetString(&conf.ClientID, envDiscordClientID)
}

func LoadTOML(paths ...string) (config Config, err error) {
	config = DefaultConfig()

	for _, path := range paths {
		var (
			file      *os.File
			fileBytes []byte
		)

		file, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if os.IsNotExist(err) || file == nil {
			err = errors.Wrap(err, fmt.Sprintf("Unable to load config from %s", path))
			return
		}

		if fileBytes, err = io.ReadAll(file); err == nil {
			err = toml.Unmarshal(fileBytes, &config)
			if err != nil {
				return
			}
		}
	}

	return
}
