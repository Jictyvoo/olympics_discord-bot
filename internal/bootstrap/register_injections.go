package bootstrap

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra"
)

func DoInjections(inj remy.Injector) {
	remy.RegisterInstance(inj, ".rest_cache", "cacheDirectory")
	infra.RegisterInfraServices(inj)
	domain.RegisterUCs(inj)
	domain.RegisterServices(inj)
}
