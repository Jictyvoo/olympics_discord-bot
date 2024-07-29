package config

import "strconv"

const DefaultFileName = "conf.toml"

type (
	Server struct {
		Host string `toml:"address"`
		Port uint16 `toml:"port"`
	}

	Discord struct {
		Token    string `toml:"token"`
		ClientID string `toml:"client_id"`
	}

	Config struct {
		IsDebug     bool    `toml:"is_debug"`
		ProjectName string  `toml:"project_name"`
		Server      Server  `toml:"server"`
		Discord     Discord `toml:"discord"`
	}
)

func (conf Server) Address() string {
	return conf.Host + ":" + strconv.Itoa(int(conf.Port))
}
