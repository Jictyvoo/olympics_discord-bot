package notifier

//go:generate go tool -modfile=../../../../tools/go.mod mockgen -source=interfaces.go -destination=interfaces_mock_test.go -package=notifier
//go:generate go tool -modfile=../../../../tools/go.mod mockgen -source=render/render.go -destination=render_mock_test.go -package=notifier

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// FixtureReader lists fixtures eligible for notification evaluation.
type FixtureReader interface {
	ListFixturesStartingBefore(before time.Time) ([]eventcore.Fixture, error)
}

// NotificationRepo persists and queries notification records.
type NotificationRepo interface {
	GetNotificationByChecksum(checksum string) (eventcore.Notification, error)
	UpsertNotification(n eventcore.Notification) error
	UpdateNotificationStatus(
		id eventcore.CanonicalID,
		status eventcore.NotificationStatus,
	) error
}

// Dispatcher sends a notification message to a channel.
type Dispatcher interface {
	Send(channelID, content string) (messageID string, err error)
}

// MentionResolver resolves the user IDs to @mention for a fixture's facts.
type MentionResolver interface {
	MentionsFor(
		guildID string, countryCodes []string, disciplineCode string,
	) ([]string, error)
}

// ResultReader loads the results recorded for a fixture.
type ResultReader interface {
	ListResultsByFixture(fixtureID eventcore.CanonicalID) ([]eventcore.Result, error)
}

// CompetitionReader resolves the competition that owns a fixture.
type CompetitionReader interface {
	GetCompetitionByFixture(fixtureID eventcore.CanonicalID) (eventcore.Competition, error)
}

// ParticipantReader loads the participants taking part in a fixture.
type ParticipantReader interface {
	ListParticipantsByFixture(fixtureID eventcore.CanonicalID) ([]eventcore.Participant, error)
}
