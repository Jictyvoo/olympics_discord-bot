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

// DoInjections is the only place that switches on conf.Database.Driver and
// conf.Providers[].Code.
func DoInjections(inj remy.Injector, conf appconfig.Config, db *sql.DB) {
	remy.RegisterInstance(inj, conf)

	// Fallback context for startup-time resolutions with no per-op ctx;
	// remy.GetWithContext overrides it per call.
	remy.RegisterInstance[context.Context](inj, context.Background())

	httpdatasource.Register(inj)
	registerCache(inj, conf)

	registerRepositories(inj, conf, db)
	registerProviders(inj, conf)
	subscriptions.Register(inj)

	// Subject the syncer emits fixtures on; observers subscribe in WireObservers.
	remy.RegisterInstance(inj, services.NewSubject[eventcore.Fixture]())

	syncer.Register(inj, conf.Runtime.SyncInterval)
	notifier.Register(
		inj,
		conf.Discord.DefaultChannel,
		providerChannels(conf),
		conf.Discord.GuildID,
		conf.Runtime.NotifyWindow,
	)
	discordsync.Register(inj, conf.Discord.GuildID, conf.Runtime.DiscordHorizon)
}

// providerChannels maps each enabled provider that declares a discord_channel to
// its channel name; providers without one fall back to the default channel.
func providerChannels(conf appconfig.Config) map[eventcore.ProviderID]string {
	channels := make(map[eventcore.ProviderID]string, len(conf.Providers))
	for _, pc := range conf.Providers {
		if pc.Enabled && pc.DiscordChannel != "" {
			channels[pc.Code] = pc.DiscordChannel
		}
	}
	return channels
}
