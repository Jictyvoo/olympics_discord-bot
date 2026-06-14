package eventcore

import "time"

type AlertKind string

const (
	AlertFixtureStart AlertKind = "fixture_start"
	AlertFixtureEnd   AlertKind = "fixture_end"
	AlertResultUpdate AlertKind = "result_update"
	AlertStatusChange AlertKind = "status_change"
)

type Alert struct {
	ID        CanonicalID
	FixtureID CanonicalID
	Kind      AlertKind
	Payload   map[string]any
}

type NotificationStatus string

const (
	NotificationPending   NotificationStatus = "pending"
	NotificationSent      NotificationStatus = "sent"
	NotificationFailed    NotificationStatus = "failed"
	NotificationCancelled NotificationStatus = "cancelled"
	NotificationSkipped   NotificationStatus = "skipped"
)

func (s NotificationStatus) Valid() bool {
	switch s {
	case NotificationPending, NotificationSent, NotificationFailed,
		NotificationCancelled, NotificationSkipped:
		return true
	}
	return false
}

// ParseNotificationStatus defaults unknown values to NotificationPending.
func ParseNotificationStatus(s string) NotificationStatus {
	ns := NotificationStatus(s)
	if ns.Valid() {
		return ns
	}
	return NotificationPending
}

type Notification struct {
	ID        CanonicalID
	AlertID   CanonicalID
	ChannelID string
	MessageID string
	Status    NotificationStatus
	Checksum  string
	SentAt    time.Time
}
