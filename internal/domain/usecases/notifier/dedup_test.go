package notifier

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// A prior non-sent record (e.g. a previous Failed attempt) must not short-circuit
// dedup: the fixture should be re-dispatched and marked sent.
func TestNotifier_ChecksumChanged_ReNotifies(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("renotify")

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("renotify").Return(eventcore.Notification{
		AlertID: f.ID, Status: eventcore.NotificationFailed, Checksum: "renotify",
	}, nil)
	var upserted []eventcore.Notification
	repo.EXPECT().UpsertNotification(gomock.Any()).
		Do(func(nt eventcore.Notification) { upserted = append(upserted, nt) }).
		Return(nil).Times(2)

	disp := NewMockDispatcher(ctrl)
	disp.EXPECT().Send(defaultChan, gomock.Any()).Return("m1", nil).Times(1)
	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().MentionsFor(testGuild, gomock.Any(), gomock.Any()).Return(nil, nil)

	n := newTestNotifier(ctrl, NewMockFixtureReader(ctrl), repo, disp, mentions, defaultChan)
	n.On(f)

	if len(upserted) != 2 || upserted[1].Status != eventcore.NotificationSent {
		t.Fatalf("expected pending -> sent re-notify, got %+v", upserted)
	}
}

// An out-of-window fixture that already has a prior record is cancelled (not
// merely skipped) and never dispatched.
func TestNotifier_OutOfWindow_WithPriorRecord_RecordsCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := eventcore.Fixture{
		ID:       eventcore.NewID(eventcore.ProviderOlympics, "cancel"),
		StartsAt: time.Now().Add(-20 * time.Hour),
		EndsAt:   time.Now().Add(-19 * time.Hour),
		Status:   eventcore.FixtureFinished,
		Checksum: "cancel",
	}

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("cancel").Return(eventcore.Notification{
		Status: eventcore.NotificationFailed, Checksum: "cancel",
	}, nil)
	var got eventcore.Notification
	repo.EXPECT().UpsertNotification(gomock.Any()).
		Do(func(nt eventcore.Notification) { got = nt }).Return(nil)

	disp := NewMockDispatcher(ctrl) // no Send expected
	n := newTestNotifier(ctrl, NewMockFixtureReader(ctrl), repo, disp,
		NewMockMentionResolver(ctrl), defaultChan)

	n.On(f)
	if got.Status != eventcore.NotificationCancelled {
		t.Fatalf("want cancelled, got %s", got.Status)
	}
}

// A future fixture outside the window simply waits: no record, no dispatch.
func TestNotifier_FutureFixtureOutOfWindow_NoRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := eventcore.Fixture{
		ID:       eventcore.NewID(eventcore.ProviderOlympics, "future"),
		StartsAt: time.Now().Add(50 * time.Hour),
		EndsAt:   time.Now().Add(52 * time.Hour),
		Status:   eventcore.FixtureScheduled,
		Checksum: "future",
	}

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("future").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	// No UpsertNotification, no Send.

	disp := NewMockDispatcher(ctrl)
	n := newTestNotifier(ctrl, NewMockFixtureReader(ctrl), repo, disp,
		NewMockMentionResolver(ctrl), defaultChan)

	n.On(f)
}

// A dedup lookup failure that is not ErrNoRows must propagate, not dispatch.
func TestNotifier_DedupLookupError_Propagates(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("lookuperr")
	want := errors.New("db down")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)
	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("lookuperr").
		Return(eventcore.Notification{}, want)

	n := newTestNotifier(ctrl, fixtures, repo, NewMockDispatcher(ctrl),
		NewMockMentionResolver(ctrl), defaultChan)

	if err := n.NotifyPending(); !errors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

func TestNotifier_OutOfWindow_RecordsSkipped(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := eventcore.Fixture{
		ID:       eventcore.NewID(eventcore.ProviderOlympics, "old"),
		StartsAt: time.Now().Add(-20 * time.Hour),
		EndsAt:   time.Now().Add(-19 * time.Hour),
		Status:   eventcore.FixtureFinished,
		Checksum: "old",
	}

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("old").Return(eventcore.Notification{}, sql.ErrNoRows)
	var got eventcore.Notification
	repo.EXPECT().
		UpsertNotification(gomock.Any()).
		Do(func(n eventcore.Notification) { got = n }).
		Return(nil)

	disp := NewMockDispatcher(ctrl) // no Send expected
	n := newTestNotifier(
		ctrl,
		NewMockFixtureReader(ctrl),
		repo,
		disp,
		NewMockMentionResolver(ctrl),
		defaultChan,
	)

	n.On(f)
	if got.Status != eventcore.NotificationSkipped {
		t.Fatalf("want skipped, got %s", got.Status)
	}
}

func TestNotifier_WithinWindow_Boundary(t *testing.T) {
	n := &Notifier{window: time.Hour}
	now := time.Now()

	// startDiff+endDiff == 0 <= 2*window -> eligible.
	inside := eventcore.Fixture{StartsAt: now, EndsAt: now}
	if !n.withinWindow(inside) {
		t.Error("fixture at now should be within window")
	}
	// startDiff+endDiff == 4h > 2*window (2h) -> ineligible.
	outside := eventcore.Fixture{StartsAt: now.Add(2 * time.Hour), EndsAt: now.Add(2 * time.Hour)}
	if n.withinWindow(outside) {
		t.Error("fixture 2h out should be outside window")
	}
}
