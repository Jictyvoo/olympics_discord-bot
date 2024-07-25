package domain

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
)

func RegisterUCs(inj remy.Injector) {
	remy.RegisterConstructorArgs2(
		inj,
		remy.Factory[usecases.FetcherCacheUseCase],
		usecases.NewFetcherCacheUseCase,
	)
}
