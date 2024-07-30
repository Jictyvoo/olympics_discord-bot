package bootstrap

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra"
	"github.com/jictyvoo/olympics_data_fetcher/pkg/config"
)

func DoInjections(inj remy.Injector, conf config.Config) {
	remy.RegisterInstance(inj, conf)
	remy.RegisterInstance(inj, entities.GetLanguage(conf.Runtime.APILocale))
	remy.RegisterInstance(inj, ".rest_cache", "cacheDirectory")
	infra.RegisterInfraServices(inj)
	domain.RegisterUCs(inj)
	domain.RegisterServices(inj)
}
