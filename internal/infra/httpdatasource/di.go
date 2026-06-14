package httpdatasource

import "github.com/wrapped-owls/goremy-di/remy"

func Register(inj remy.Injector) {
	remy.RegisterConstructor(inj, remy.Factory[Client], func() Client {
		return New()
	})
}
