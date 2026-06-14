package eventcore

import "time"

type DiscordEventStatus string

const (
	DiscordEventScheduled DiscordEventStatus = "scheduled"
	DiscordEventActive    DiscordEventStatus = "active"
	DiscordEventCompleted DiscordEventStatus = "completed"
	DiscordEventCancelled DiscordEventStatus = "cancelled"
)

func (s DiscordEventStatus) Valid() bool {
	switch s {
	case DiscordEventScheduled, DiscordEventActive, DiscordEventCompleted, DiscordEventCancelled:
		return true
	}
	return false
}

type DiscordEvent struct {
	FixtureID      CanonicalID
	GuildID        string
	DiscordEventID string
	Status         DiscordEventStatus
	LastChecksum   string
	UpdatedAt      time.Time
}
