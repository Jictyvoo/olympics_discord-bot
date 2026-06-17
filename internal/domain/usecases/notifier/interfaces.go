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
	GetLatestSentNotificationByAlert(
		alertID eventcore.CanonicalID,
	) (eventcore.Notification, error)
	UpsertNotification(n eventcore.Notification) error
	UpdateNotificationStatus(
		id eventcore.CanonicalID,
		status eventcore.NotificationStatus,
	) error
}

// Dispatcher sends a notification message to a channel, or edits one in place.
type Dispatcher interface {
	Send(channelID, content string) (messageID string, err error)
	Edit(channelID, messageID, content string) error
}

// ChannelEnsurer resolves a channel name to its ID, creating the channel if it
// does not already exist.
type ChannelEnsurer interface {
	ResolveChannel(guildID, channelName string) (string, error)
}

// MentionResolver resolves the user IDs to @mention for a fixture's facts.
type MentionResolver interface {
	MentionsFor(
		guildID string, countryCodes []string, disciplineCode string,
	) ([]string, error)
}

// FixtureContextReader resolves the competition, stage and group locating a
// fixture in a single query.
type FixtureContextReader interface {
	GetFixtureContext(fixtureID eventcore.CanonicalID) (eventcore.FixtureContext, error)
}

// CompetitorReader loads each competitor of a fixture together with its role,
// resolved flag code and result.
type CompetitorReader interface {
	ListFixtureCompetitors(
		fixtureID eventcore.CanonicalID,
	) ([]eventcore.FixtureCompetitor, error)
}
