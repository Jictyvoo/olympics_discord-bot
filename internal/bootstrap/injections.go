package bootstrap

import (
	"context"
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/discordsync"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/subscriptions"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/syncer"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

// DoInjections wires the full application into the DI container.
// This is the only place that switches on conf.Database.Driver and conf.Providers[].Code.
func DoInjections(inj remy.Injector, conf appconfig.Config, db *sql.DB) {
	remy.RegisterInstance(inj, conf)
	remy.RegisterInstance(inj, db)

	// Fallback context for startup-time resolutions that have no per-op ctx
	// (observers, serve.go, provider/di.go). remy.GetWithContext overrides this
	// per call via its sub-injector. Registered under the context.Context
	// interface type so factories requesting a context.Context arg resolve it.
	remy.RegisterInstance[context.Context](inj, context.Background())

	httpdatasource.Register(inj)
	registerCache(inj, conf)

	registerRepositories(inj, conf, db)
	registerProviders(inj, conf)
	subscriptions.Register(inj)

	// Long-lived subject the syncer emits persisted fixtures on; observers
	// (notifier, discordsync) subscribe in WireObservers.
	remy.RegisterInstance(inj, services.NewSubject[eventcore.Fixture]())

	syncer.Register(inj, conf.Runtime.SyncInterval)
	notifier.Register(
		inj,
		conf.Discord.DefaultChannel,
		conf.Discord.GuildID,
		conf.Runtime.NotifyWindow,
	)
	discordsync.Register(inj, conf.Discord.GuildID, conf.Runtime.DiscordHorizon)
}
