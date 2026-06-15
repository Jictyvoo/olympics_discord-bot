package discordsync

import (
	"time"

	"github.com/wrapped-owls/goremy-di/remy"
)

func Register(inj remy.Injector, guildID string, horizon time.Duration) {
	remy.RegisterConstructorArgs3(
		inj,
		remy.Factory[*DiscordSync],
		func(fr FixtureReader, dr DiscordEventRepo, facade ScheduledEventFacade) *DiscordSync {
			return New(fr, dr, facade, guildID, horizon)
		},
	)
}
