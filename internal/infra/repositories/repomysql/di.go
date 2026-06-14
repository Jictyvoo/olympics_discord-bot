package repomysql

import (
	"database/sql"

	"github.com/wrapped-owls/goremy-di/remy"
)

// Register wires the base and every repo as factories. The base reads the
// injected context, so remy.GetWithContext flows the per-op ctx to its timeouts.
func Register(inj remy.Injector, db *sql.DB) {
	remy.RegisterInstance(inj, db)

	remy.RegisterConstructorArgs2(inj, remy.Factory[*repoMySQL], newRepo)

	remy.RegisterConstructorArgs1(inj, remy.Factory[*Repository], NewRepository)
	remy.RegisterConstructorArgs1(inj, remy.Factory[FixtureRepo], NewFixtureRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[CompetitionRepo], NewCompetitionRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[SeasonRepo], NewSeasonRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[StageRepo], NewStageRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[GroupRepo], NewGroupRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[VenueRepo], NewVenueRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[StandingRepo], NewStandingRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[ParticipantRepo], NewParticipantRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[ResultRepo], NewResultRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[SyncStateRepo], NewSyncStateRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[NotificationRepo], NewNotificationRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[DiscordEventRepo], NewDiscordEventRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[CountryRepo], NewCountryRepo)
	remy.RegisterConstructorArgs1(inj, remy.Factory[SubscriptionRepo], NewSubscriptionRepo)
}
