package notifier

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/pkg/idgen"
)

// withinWindow reports whether the fixture's start and end both fall close
// enough to now to be worth notifying.
func (n *Notifier) withinWindow(f eventcore.Fixture) bool {
	now := time.Now()
	startDiff := absDuration(f.StartsAt.Sub(now))
	endDiff := absDuration(f.EndsAt.Sub(now))
	return startDiff+endDiff <= 2*n.window
}

// recordIneligible persists a terminal status for a fixture that fell outside
// the window. Future ones simply wait; only well-past ones are recorded as
// skipped (no prior record) or cancelled (a prior record exists).
func (n *Notifier) recordIneligible(
	f eventcore.Fixture,
	checksum string,
	existing eventcore.Notification,
	channelID string,
) error {
	if !f.EndsAt.Before(time.Now().Add(-n.window / 2)) { //nolint:mnd // half the notify window
		return nil
	}
	status := eventcore.NotificationSkipped
	if existing.Status != "" {
		status = eventcore.NotificationCancelled
	}
	return n.notifs.UpsertNotification(
		eventcore.Notification{
			ID:        idgen.NewV7(),
			AlertID:   f.ID,
			ChannelID: channelID,
			Status:    status,
			Checksum:  checksum,
		},
	)
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
