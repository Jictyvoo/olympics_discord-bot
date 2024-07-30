package domain

import (
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/services"
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
)

func RegisterUCs(inj remy.Injector) {
	remy.RegisterConstructorArgs2(
		inj,
		remy.Factory[usecases.FetcherCacheUseCase],
		usecases.NewFetcherCacheUseCase,
	)
}

func RegisterServices(inj remy.Injector) {
	remy.RegisterConstructorArgs2Err(
		inj, remy.Factory[services.EventObserver],
		services.NewOlympicEventManager,
	)

	remy.Register(
		inj, remy.LazySingleton(
			func(retriever remy.DependencyRetriever) (*services.EventNotifier, error) {
				cancelChan, _ := remy.DoGet[services.CancelChannel](retriever)
				if cancelChan == nil {
					cancelChan = make(services.CancelChannel, 1)
				}

				fetcherUC, err := remy.DoGet[usecases.FetcherCacheUseCase](retriever)
				if err != nil {
					return nil, err
				}

				var loaderRepo services.EventNotifierRepository
				if loaderRepo, err = remy.DoGet[services.EventNotifierRepository](retriever); err != nil {
					return nil, err
				}
				return services.NewEventNotifier(cancelChan, 4*time.Minute, loaderRepo, fetcherUC)
			},
		),
	)
}
