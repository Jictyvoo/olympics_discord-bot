package bootstrap

import (
	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore"
)

func registerCache(inj remy.Injector, conf appconfig.Config) {
	switch conf.Cache.Backend {
	case "memory":
		cachestore.RegisterMemory(inj, conf.Cache.TTL)
	default:
		cachestore.RegisterFile(inj, conf.Cache.FilePath, conf.Cache.TTL)
	}
}
