package notifier

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
	"github.com/jictyvoo/olhojogo/pkg/idgen"
)

const (
	defaultWindowHours   = 3
	defaultWindowMinutes = 30
)

// Notifier evaluates pending fixtures and dispatches notifications.
type Notifier struct {
	fixtures    FixtureReader
	notifs      NotificationRepo
	dispatch    Dispatcher
	context     FixtureContextReader
	competitors CompetitorReader
	renderer    render.Renderer
	mentions    MentionResolver
	channelID   string
	guildID     string
	window      time.Duration // fixtures whose start/end straddle this window are eligible
}

func New(
	fixtures FixtureReader,
	notifs NotificationRepo,
	dispatch Dispatcher,
	context FixtureContextReader,
	competitors CompetitorReader,
	renderer render.Renderer,
	mentions MentionResolver,
	channelID, guildID string,
	window time.Duration,
) *Notifier {
	if window == 0 {
		window = defaultWindowHours*time.Hour + defaultWindowMinutes*time.Minute
	}
	return &Notifier{
		fixtures:    fixtures,
		notifs:      notifs,
		dispatch:    dispatch,
		context:     context,
		competitors: competitors,
		renderer:    renderer,
		mentions:    mentions,
		channelID:   channelID,
		guildID:     guildID,
		window:      window,
	}
}

// On satisfies services.Observer[eventcore.Fixture]: it is invoked for every
// fixture the syncer persists, and decides whether to notify.
func (n *Notifier) On(f eventcore.Fixture) {
	if err := n.evaluate(f); err != nil {
		slog.Error(
			"notifier: evaluate",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	}
}

// NotifyPending evaluates all fixtures starting before now+window. It is an
// optional startup catch-up sweep; steady-state notification is event-driven
// through On.
func (n *Notifier) NotifyPending() error {
	horizon := time.Now().Add(n.window)
	candidates, err := n.fixtures.ListFixturesStartingBefore(horizon)
	if err != nil {
		return err
	}

	var errs []error
	for _, f := range candidates {
		if err = n.evaluate(f); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (n *Notifier) evaluate(f eventcore.Fixture) error {
	if n.channelID == "" {
		return nil
	}

	// The fixture checksum covers status, times and results, so it changes
	// whenever the rendered content would.
	checksum := f.Checksum

	// Edit the existing message in place when content changed, not gated on the
	// window so the final result always lands.
	prior, err := n.notifs.GetLatestSentNotificationByAlert(f.ID)
	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return err
	case err == nil && prior.MessageID != "":
		if prior.Checksum == checksum {
			return nil
		}
		return n.edit(f, prior, checksum)
	}

	// No message yet: dedup on the checksum and gate on the window.
	existing, err := n.notifs.GetNotificationByChecksum(checksum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existing.Status == eventcore.NotificationSent {
		return nil
	}
	if !n.withinWindow(f) {
		return n.recordIneligible(f, checksum, existing)
	}

	content, err := n.compose(f)
	if err != nil {
		return err
	}
	return n.send(f, checksum, content)
}

// edit updates an already-sent message in place and records the new checksum,
// so a fixture's notification evolves from scheduled to finished as one message.
func (n *Notifier) edit(
	f eventcore.Fixture,
	prior eventcore.Notification,
	checksum string,
) error {
	content, err := n.compose(f)
	if err != nil {
		return err
	}
	if err = n.dispatch.Edit(n.channelID, prior.MessageID, content); err != nil {
		slog.Error(
			"notifier: edit failed",
			slog.String("fixture", f.ID.String()),
			slog.String("message", prior.MessageID),
			slog.String("err", err.Error()),
		)
		return err
	}
	prior.Checksum = checksum
	prior.Status = eventcore.NotificationSent
	prior.SentAt = time.Now()
	if err = n.notifs.UpsertNotification(prior); err != nil {
		slog.Error("notifier: update edited notification", slog.String("err", err.Error()))
	}
	return nil
}

func (n *Notifier) send(f eventcore.Fixture, checksum, content string) error {
	// Write a pending row before dispatching so a crash doesn't cause infinite retries.
	notification := eventcore.Notification{
		ID:        idgen.NewV7(),
		AlertID:   f.ID,
		ChannelID: n.channelID,
		Status:    eventcore.NotificationPending,
		Checksum:  checksum,
	}
	if err := n.notifs.UpsertNotification(notification); err != nil {
		return err
	}

	msgID, sendErr := n.dispatch.Send(n.channelID, content)
	if sendErr != nil {
		slog.Error(
			"notifier: send failed",
			slog.String("fixture", f.ID.String()),
			slog.String("channel", n.channelID),
			slog.String("err", sendErr.Error()),
		)
		_ = n.notifs.UpdateNotificationStatus(notification.ID, eventcore.NotificationFailed)
		return sendErr
	}

	notification.MessageID = msgID
	notification.Status = eventcore.NotificationSent
	notification.SentAt = time.Now()
	if err := n.notifs.UpsertNotification(notification); err != nil {
		slog.Error("notifier: update sent status", slog.String("err", err.Error()))
	}
	return nil
}
