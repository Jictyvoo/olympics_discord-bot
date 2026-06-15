package notifier

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
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
	fixtures     FixtureReader
	notifs       NotificationRepo
	dispatch     Dispatcher
	results      ResultReader
	competitions CompetitionReader
	participants ParticipantReader
	renderer     render.Renderer
	mentions     MentionResolver
	channelID    string
	guildID      string
	window       time.Duration // fixtures whose start/end straddle this window are eligible
}

func New(
	fixtures FixtureReader,
	notifs NotificationRepo,
	dispatch Dispatcher,
	results ResultReader,
	competitions CompetitionReader,
	participants ParticipantReader,
	renderer render.Renderer,
	mentions MentionResolver,
	channelID, guildID string,
	window time.Duration,
) *Notifier {
	if window == 0 {
		window = defaultWindowHours*time.Hour + defaultWindowMinutes*time.Minute
	}
	return &Notifier{
		fixtures:     fixtures,
		notifs:       notifs,
		dispatch:     dispatch,
		results:      results,
		competitions: competitions,
		participants: participants,
		renderer:     renderer,
		mentions:     mentions,
		channelID:    channelID,
		guildID:      guildID,
		window:       window,
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
		if err := n.evaluate(f); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (n *Notifier) evaluate(f eventcore.Fixture) error {
	if n.channelID == "" {
		return nil
	}

	// While the fixture is still ongoing, dedup on a checksum that ignores
	// results so mid-event result updates don't trigger repeat notifications;
	// once it finishes, the full checksum (results included) governs dedup.
	checksum := f.Checksum
	if f.Status != eventcore.FixtureFinished &&
		time.Now().Before(f.EndsAt.Add(n.window/2)) {
		checksum = f.ComputeChecksum()
	}

	existing, err := n.notifs.GetNotificationByChecksum(checksum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existing.Status == eventcore.NotificationSent {
		return nil
	}

	// Time-window gating: a fixture is eligible only while its start/end straddle
	// the allowed window. A fixture that is already well past is recorded as
	// skipped/cancelled rather than sent.
	if !n.withinWindow(f) {
		return n.recordIneligible(f, checksum, existing)
	}

	view := n.buildView(f)
	content := n.renderer.Render(view)
	if prefix, mErr := n.mentionPrefix(view); mErr != nil {
		return mErr
	} else if prefix != "" {
		content = prefix + content
	}

	return n.send(f, checksum, content)
}

// withinWindow reports whether the fixture's start and end both fall close
// enough to now to be worth notifying (mirrors the legacy allowed-time-diff).
func (n *Notifier) withinWindow(f eventcore.Fixture) bool {
	now := time.Now()
	startDiff := absDuration(f.StartsAt.Sub(now))
	endDiff := absDuration(f.EndsAt.Sub(now))
	return startDiff+endDiff <= 2*n.window
}

// recordIneligible persists a terminal status for a fixture that fell outside
// the window. Future fixtures simply wait; only well-past ones are recorded as
// skipped (no prior record) or cancelled (a prior record exists).
func (n *Notifier) recordIneligible(
	f eventcore.Fixture,
	checksum string,
	existing eventcore.Notification,
) error {
	if !f.EndsAt.Before(time.Now().Add(-n.window / 2)) {
		return nil
	}
	status := eventcore.NotificationSkipped
	if existing.Status != "" {
		status = eventcore.NotificationCancelled
	}
	return n.notifs.UpsertNotification(eventcore.Notification{
		ID:        idgen.NewV7(),
		AlertID:   f.ID,
		ChannelID: n.channelID,
		Status:    status,
		Checksum:  checksum,
	})
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// mentionPrefix resolves the subscribers to @mention for the fixture's facts and
// renders them as a "<@id> <@id> " prefix, or "" when there are no mentions.
func (n *Notifier) mentionPrefix(
	view render.FixtureView,
) (string, error) {
	if n.mentions == nil {
		return "", nil
	}
	countryCodes := make([]string, 0, len(view.Participants))
	for _, p := range view.Participants {
		if p.CountryISO != "" {
			countryCodes = append(countryCodes, p.CountryISO)
		}
	}
	users, err := n.mentions.MentionsFor(n.guildID, countryCodes, view.Competition.Code)
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", nil
	}
	var b strings.Builder
	for _, id := range users {
		fmt.Fprintf(&b, "<@%s> ", id)
	}
	return b.String(), nil
}

// send dispatches the fixture notification to the configured channel, deduping
// per checksum and recording the notification lifecycle.
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

// buildView assembles the render aggregate for a fixture. Enrichment reads are
// best-effort: a failure is logged and the view degrades gracefully rather than
// blocking the notification, whose essential payload is the fixture itself.
func (n *Notifier) buildView(f eventcore.Fixture) render.FixtureView {
	view := render.FixtureView{Fixture: f}

	results, err := n.results.ListResultsByFixture(f.ID)
	if err != nil {
		slog.Warn(
			"notifier: load results",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	} else {
		view.Results = results
	}

	competition, err := n.competitions.GetCompetitionByFixture(f.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Warn(
			"notifier: load competition",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	} else {
		view.Competition = competition
	}

	participants, err := n.participants.ListParticipantsByFixture(f.ID)
	if err != nil {
		slog.Warn(
			"notifier: load participants",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	} else {
		view.Participants = participants
	}

	return view
}
