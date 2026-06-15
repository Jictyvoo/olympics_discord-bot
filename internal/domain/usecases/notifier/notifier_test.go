package notifier

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
)

const (
	defaultChan = "chan"
	testGuild   = "g1"
)

// mkFixture builds a just-finished fixture inside the notify window. Being
// finished, dedup uses its literal Checksum (not the results-free one), which
// keeps these tests asserting against the given checksum string.
func mkFixture(checksum string) eventcore.Fixture {
	return eventcore.Fixture{
		ID:       eventcore.NewID(eventcore.ProviderOlympics, "evt-"+checksum),
		Name:     "Match " + checksum,
		StartsAt: time.Now().Add(-2 * time.Hour),
		EndsAt:   time.Now().Add(-time.Hour),
		Status:   eventcore.FixtureFinished,
		Checksum: checksum,
	}
}

// newTestNotifier builds a Notifier with stub enrichment readers and the given
// collaborators. Enrichment reads are wired with AnyTimes returns by default.
func newTestNotifier(
	ctrl *gomock.Controller,
	fixtures FixtureReader,
	repo NotificationRepo,
	disp Dispatcher,
	mentions MentionResolver,
	channelID string,
) *Notifier {
	resultsReader := NewMockResultReader(ctrl)
	resultsReader.EXPECT().ListResultsByFixture(gomock.Any()).
		Return(nil, nil).AnyTimes()
	comps := NewMockCompetitionReader(ctrl)
	comps.EXPECT().GetCompetitionByFixture(gomock.Any()).
		Return(eventcore.Competition{}, nil).AnyTimes()
	participants := NewMockParticipantReader(ctrl)
	participants.EXPECT().ListParticipantsByFixture(gomock.Any()).
		Return(nil, nil).AnyTimes()
	rndr := NewMockRenderer(ctrl)
	rndr.EXPECT().Render(gomock.Any()).Return("x").AnyTimes()

	return New(
		fixtures, repo, disp, resultsReader, comps, participants, rndr,
		mentions, channelID, testGuild, 0,
	)
}

// enrichedReaders builds the three enrichment readers, each expecting one call
// returning the given values.
func enrichedReaders(
	ctrl *gomock.Controller,
	results []eventcore.Result,
	comp eventcore.Competition,
	parts []eventcore.Participant,
) (*MockResultReader, *MockCompetitionReader, *MockParticipantReader) {
	resultsReader := NewMockResultReader(ctrl)
	resultsReader.EXPECT().ListResultsByFixture(gomock.Any()).Return(results, nil)
	comps := NewMockCompetitionReader(ctrl)
	comps.EXPECT().GetCompetitionByFixture(gomock.Any()).Return(comp, nil)
	participants := NewMockParticipantReader(ctrl)
	participants.EXPECT().ListParticipantsByFixture(gomock.Any()).Return(parts, nil)
	return resultsReader, comps, participants
}

func TestNotifier_NotifyPending_SkipsAlreadySent(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("abc")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetNotificationByChecksum("abc").
		Return(eventcore.Notification{
			AlertID: f.ID, Status: eventcore.NotificationSent, Checksum: "abc",
		}, nil)
	// No UpsertNotification expected for an already-sent fixture.

	disp := NewMockDispatcher(ctrl)
	mentions := NewMockMentionResolver(ctrl)

	n := newTestNotifier(ctrl, fixtures, repo, disp, mentions, defaultChan)

	if err := n.NotifyPending(); err != nil {
		t.Fatalf("NotifyPending: %v", err)
	}
}

func TestNotifier_NotifyPending_DispatchesAndMarksSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("def")
	results := []eventcore.Result{{FixtureID: f.ID, Outcome: eventcore.OutcomeMedalGold}}
	comp := eventcore.Competition{Code: "ATH", Discipline: "Athletics"}
	parts := []eventcore.Participant{{Name: "Usain Bolt", CountryISO: "JAM"}}

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetNotificationByChecksum("def").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	var upserted []eventcore.Notification
	repo.EXPECT().
		UpsertNotification(gomock.Any()).
		Do(func(nt eventcore.Notification) { upserted = append(upserted, nt) }).
		Return(nil).
		Times(2)

	var sentContent string
	disp := NewMockDispatcher(ctrl)
	disp.EXPECT().
		Send(defaultChan, gomock.Any()).
		Do(func(_, content string) { sentContent = content }).
		Return("msg-42", nil).
		Times(1)

	resultsReader, comps, participants := enrichedReaders(ctrl, results, comp, parts)
	rndr := NewMockRenderer(ctrl)
	var lastView render.FixtureView
	rndr.EXPECT().
		Render(gomock.Any()).
		Do(func(view render.FixtureView) { lastView = view }).
		Return("x")

	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().
		MentionsFor(testGuild, []string{"JAM"}, "ATH").
		Return([]string{"u1", "u2"}, nil)

	n := New(
		fixtures, repo, disp, resultsReader, comps, participants, rndr,
		mentions, defaultChan, testGuild, 0,
	)

	if err := n.NotifyPending(); err != nil {
		t.Fatalf("NotifyPending: %v", err)
	}
	if sentContent != "<@u1> <@u2> x" {
		t.Fatalf("content = %q, want mention prefix + body", sentContent)
	}
	assertDispatchedAndSent(t, lastView, upserted)
}

func assertDispatchedAndSent(
	t *testing.T, lastView render.FixtureView, upserted []eventcore.Notification,
) {
	t.Helper()
	if len(lastView.Results) != 1 ||
		lastView.Competition.Code != "ATH" || len(lastView.Participants) != 1 {
		t.Fatalf("renderer did not receive enriched view: %+v", lastView)
	}
	if len(upserted) != 2 {
		t.Fatalf("expected 2 upserts (pending -> sent); got %d", len(upserted))
	}
	if got := upserted[0].Status; got != eventcore.NotificationPending {
		t.Fatalf("first upsert status: want pending, got %s", got)
	}
	final := upserted[1]
	if final.Status != eventcore.NotificationSent {
		t.Fatalf("second upsert status: want sent, got %s", final.Status)
	}
	if final.MessageID != "msg-42" || final.ChannelID != defaultChan {
		t.Fatalf("final notification message/channel mismatch: %+v", final)
	}
}

func TestNotifier_NotifyPending_NoMentions_NoPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("nopfx")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetNotificationByChecksum("nopfx").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	repo.EXPECT().UpsertNotification(gomock.Any()).Return(nil).AnyTimes()

	var sentContent string
	disp := NewMockDispatcher(ctrl)
	disp.EXPECT().
		Send(defaultChan, gomock.Any()).
		Do(func(_, content string) { sentContent = content }).
		Return("m", nil).
		Times(1)

	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().
		MentionsFor(testGuild, gomock.Any(), gomock.Any()).
		Return(nil, nil)

	n := newTestNotifier(ctrl, fixtures, repo, disp, mentions, defaultChan)

	if err := n.NotifyPending(); err != nil {
		t.Fatalf("NotifyPending: %v", err)
	}
	if sentContent != "x" {
		t.Fatalf("content = %q, want body without prefix", sentContent)
	}
}

func TestNotifier_NotifyPending_DispatchFailure_MarksFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("xyz")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetNotificationByChecksum("xyz").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	repo.EXPECT().UpsertNotification(gomock.Any()).Return(nil).AnyTimes()
	var gotStatus eventcore.NotificationStatus
	statusCalls := 0
	repo.EXPECT().
		UpdateNotificationStatus(gomock.Any(), gomock.Any()).
		Do(func(_ eventcore.CanonicalID, st eventcore.NotificationStatus) {
			gotStatus = st
			statusCalls++
		}).
		Return(nil).
		Times(1)

	disp := NewMockDispatcher(ctrl)
	disp.EXPECT().
		Send(gomock.Any(), gomock.Any()).
		Return("", errors.New("discord down"))

	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().
		MentionsFor(testGuild, gomock.Any(), gomock.Any()).
		Return(nil, nil)

	n := newTestNotifier(ctrl, fixtures, repo, disp, mentions, defaultChan)

	if err := n.NotifyPending(); err == nil {
		t.Fatal("expected dispatch error to be returned")
	}
	if statusCalls != 1 || gotStatus != eventcore.NotificationFailed {
		t.Fatalf("expected status update to failed; got status=%v calls=%d", gotStatus, statusCalls)
	}
}

func TestNotifier_NotifyPending_FixtureReaderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	want := errors.New("db blew up")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return(nil, want)

	n := New(
		fixtures,
		NewMockNotificationRepo(ctrl),
		NewMockDispatcher(ctrl),
		NewMockResultReader(ctrl),
		NewMockCompetitionReader(ctrl),
		NewMockParticipantReader(ctrl),
		NewMockRenderer(ctrl),
		NewMockMentionResolver(ctrl),
		defaultChan,
		testGuild,
		0,
	)
	if err := n.NotifyPending(); !errors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

// TestNotifier_NoChannel_Skips verifies that with no configured channel nothing
// is dispatched and no dedup lookup happens.
func TestNotifier_NoChannel_Skips(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("none")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	disp := NewMockDispatcher(ctrl)
	mentions := NewMockMentionResolver(ctrl)

	n := newTestNotifier(ctrl, fixtures, repo, disp, mentions, "")

	if err := n.NotifyPending(); err != nil {
		t.Fatalf("NotifyPending: %v", err)
	}
}

// TestNotifier_OutOfWindow_RecordsSkipped verifies a fixture that ended well
// outside the window is recorded as skipped and never dispatched.
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

// TestNotifier_Ongoing_UsesResultsFreeChecksum verifies an unfinished fixture
// is deduped on its results-free checksum, not its stored Checksum.
func TestNotifier_Ongoing_UsesResultsFreeChecksum(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := eventcore.Fixture{
		ID:       eventcore.NewID(eventcore.ProviderOlympics, "live"),
		StartsAt: time.Now().Add(-time.Hour),
		EndsAt:   time.Now().Add(time.Hour),
		Status:   eventcore.FixtureLive,
		Checksum: "stored-with-results",
	}

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetNotificationByChecksum(f.ComputeChecksum()).
		Return(eventcore.Notification{Status: eventcore.NotificationSent}, nil)
	// Already sent under the results-free checksum -> no dispatch, no upsert.

	disp := NewMockDispatcher(ctrl)
	n := newTestNotifier(
		ctrl,
		NewMockFixtureReader(ctrl),
		repo,
		disp,
		NewMockMentionResolver(ctrl),
		defaultChan,
	)

	n.On(f)
}

func TestNew_DefaultWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	n := New(
		NewMockFixtureReader(ctrl), NewMockNotificationRepo(ctrl), NewMockDispatcher(ctrl),
		NewMockResultReader(ctrl), NewMockCompetitionReader(ctrl), NewMockParticipantReader(ctrl),
		NewMockRenderer(ctrl), NewMockMentionResolver(ctrl), "c", testGuild, 0,
	)
	if n.window == 0 {
		t.Fatal("default window must be non-zero")
	}
}
