package bootstrap

import (
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/discordsync"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/subscriptions"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/syncer"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite"
)

func registerRepositories(inj remy.Injector, conf appconfig.Config, db *sql.DB) {
	switch conf.Database.Driver {
	case "mysql":
		repomysql.Register(inj, db)
		bindIface[repomysql.FixtureRepo, notifier.FixtureReader](inj)
		bindIface[repomysql.FixtureRepo, discordsync.FixtureReader](inj)
		bindIface[repomysql.CompetitionRepo, notifier.FixtureContextReader](inj)
		bindIface[repomysql.ParticipantRepo, notifier.CompetitorReader](inj)
		bindIface[repomysql.NotificationRepo, notifier.NotificationRepo](inj)
		bindIface[repomysql.DiscordEventRepo, discordsync.DiscordEventRepo](inj)
		bindIface[repomysql.VenueRepo, discordsync.VenueReader](inj)
		bindIface[repomysql.SubscriptionRepo, subscriptions.Repository](inj)
		bindIface[repomysql.CountryRepo, subscriptions.CountryLister](inj)
		bindIface[*repomysql.Repository, syncer.Repository](inj)
	default:
		reposqlite.Register(inj, db)
		bindIface[reposqlite.FixtureRepo, notifier.FixtureReader](inj)
		bindIface[reposqlite.FixtureRepo, discordsync.FixtureReader](inj)
		bindIface[reposqlite.CompetitionRepo, notifier.FixtureContextReader](inj)
		bindIface[reposqlite.ParticipantRepo, notifier.CompetitorReader](inj)
		bindIface[reposqlite.NotificationRepo, notifier.NotificationRepo](inj)
		bindIface[reposqlite.DiscordEventRepo, discordsync.DiscordEventRepo](inj)
		bindIface[reposqlite.VenueRepo, discordsync.VenueReader](inj)
		bindIface[reposqlite.SubscriptionRepo, subscriptions.Repository](inj)
		bindIface[reposqlite.CountryRepo, subscriptions.CountryLister](inj)
		bindIface[*reposqlite.Repository, syncer.Repository](inj)
	}
}

// bindIface exposes the concrete type C under the consumer interface I.
func bindIface[C any, I any](inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[I],
		func(c C) I { return any(c).(I) },
	)
}
