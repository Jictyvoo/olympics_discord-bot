package infra

import (
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/datasources"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/datasources/dsrest"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories"
)

func RegisterInfraServices(inj remy.Injector) {
	remy.RegisterSingleton(
		inj,
		func(retriever remy.DependencyRetriever) (datasources.CacheableDataSource, error) {
			cacheDir, err := remy.DoGet[string](retriever, "cacheDirectory")
			if err != nil {
				return nil, err
			}

			cacheDuration, _ := remy.DoGet[time.Duration](retriever, "cache")
			return datasources.NewDirectoryCache(cacheDir, cacheDuration)
		},
	)

	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[dsrest.RESTDataSource],
		dsrest.NewCurlCacheableDatasource,
	)
	repositories.RegisterRepositories(inj)
}
