package fifafetch

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
)

func Register(inj remy.Injector) {
	remy.RegisterConstructor(inj, remy.Factory[Provider], New)
	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[provider.Strategy],
		func(p Provider) provider.Strategy { return p },
		eventcore.ProviderFIFA,
	)
}
