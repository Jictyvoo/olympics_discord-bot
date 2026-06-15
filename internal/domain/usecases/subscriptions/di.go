package subscriptions

import "github.com/wrapped-owls/goremy-di/remy"

func Register(inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[Service],
		func(repo Repository) Service { return New(repo) },
	)
}
