package subscriptions

import "github.com/wrapped-owls/goremy-di/remy"

func Register(inj remy.Injector) {
	remy.RegisterConstructorArgs2(
		inj,
		remy.Factory[Service],
		func(repo Repository, countries CountryLister) Service { return New(repo, countries) },
	)
}
