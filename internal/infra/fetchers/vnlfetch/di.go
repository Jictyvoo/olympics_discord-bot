package vnlfetch

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

func Register(inj remy.Injector, baseURL, lang, tournaments string) {
	remy.RegisterConstructorArgs2(
		inj,
		remy.Factory[Provider],
		func(client httpdatasource.Client, cache cachestore.Cache) Provider {
			return New(client, cache, baseURL, lang, tournaments)
		},
	)
	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[provider.Strategy],
		func(p Provider) provider.Strategy { return p },
		eventcore.ProviderVNL,
	)
}
