package discordsync

//go:generate go tool -modfile=../../../../tools/go.mod mockgen -source=interfaces.go -destination=interfaces_mock_test.go -package=discordsync

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/discordfacade"
)

// FixtureReader lists upcoming fixtures for reconciliation.
type FixtureReader interface {
	ListFixturesStartingBefore(before time.Time) ([]eventcore.Fixture, error)
}

// DiscordEventRepo persists and queries the local discord_events ledger.
type DiscordEventRepo interface {
	GetDiscordEventByFixture(
		fixtureID eventcore.CanonicalID,
		guildID string,
	) (eventcore.DiscordEvent, error)
	UpsertDiscordEvent(de eventcore.DiscordEvent) error
	UpdateDiscordEventStatus(
		fixtureID eventcore.CanonicalID,
		guildID string,
		status eventcore.DiscordEventStatus,
	) error
}

// ScheduledEventFacade is the subset of discordfacade.Facade used here.
type ScheduledEventFacade interface {
	ListScheduledEvents(guildID string) ([]discordfacade.ScheduledEvent, error)
	CreateScheduledEvent(
		guildID string,
		in discordfacade.ScheduledEventInput,
	) (eventID string, err error)
	UpdateScheduledEvent(
		guildID, eventID string,
		in discordfacade.ScheduledEventInput,
	) error
	CancelScheduledEvent(guildID, eventID string) error
}
