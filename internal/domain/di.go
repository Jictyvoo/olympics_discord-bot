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

	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[usecases.CanNotifyUseCase],
		usecases.NewCanNotifyUseCase,
	)
}

func RegisterServices(inj remy.Injector) {
	remy.Register(
		inj, remy.LazySingleton(
			func(retriever remy.DependencyRetriever) (*services.EventNotifier, error) {
				cancelChan, _ := remy.DoGet[services.CancelChannel](retriever)
				if cancelChan == nil {
					cancelChan = make(services.CancelChannel, 1)
				}

				loaderRepo, err := remy.DoGet[services.EventNotifierRepository](retriever)
				if err != nil {
					return nil, err
				}

				var fetcherUC usecases.FetcherCacheUseCase
				if fetcherUC, err = remy.DoGet[usecases.FetcherCacheUseCase](retriever); err != nil {
					return nil, err
				}

				var canNotifyUC usecases.CanNotifyUseCase
				if canNotifyUC, err = remy.DoGet[usecases.CanNotifyUseCase](retriever); err != nil {
					return nil, err
				}
				return services.NewEventNotifier(
					cancelChan, 4*time.Minute,
					loaderRepo, fetcherUC, canNotifyUC,
				)
			},
		),
	)
}
