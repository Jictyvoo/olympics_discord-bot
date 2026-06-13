package entities

import "time"

type notificationStatus string

const (
	NotificationStatusPending   notificationStatus = "pending"
	NotificationStatusSent      notificationStatus = "sent"
	NotificationStatusCancelled notificationStatus = "cancelled"
	NotificationStatusFailed    notificationStatus = "failed"
	NotificationStatusSkipped   notificationStatus = "skipped"
)

//goland:noinspection GoExportedFuncWithUnexportedType
func NotificationStatus(status string) notificationStatus {
	switch notificationStatus(status) {
	case NotificationStatusPending:
		return NotificationStatusPending
	case NotificationStatusSent:
		return NotificationStatusSent
	case NotificationStatusCancelled:
		return NotificationStatusCancelled
	case NotificationStatusFailed:
		return NotificationStatusFailed
	case NotificationStatusSkipped:
		return NotificationStatusSkipped
	}

	return NotificationStatusPending
}

type Notification struct {
	ID            Identifier
	EventID       Identifier
	Status        notificationStatus
	EventChecksum string
	NotifiedAt    time.Time
}
