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

// mkFixture builds a just-finished fixture inside the notify window. Its literal
// Checksum governs dedup, so these tests assert against the given checksum string.
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

// noPriorSent wires the "no message has been posted yet" expectation so the
// notifier falls through to its first-send path.
func noPriorSent(repo *MockNotificationRepo) {
	repo.EXPECT().
		GetLatestSentNotificationByAlert(gomock.Any()).
		Return(eventcore.Notification{}, sql.ErrNoRows).
		AnyTimes()
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
	context := NewMockFixtureContextReader(ctrl)
	context.EXPECT().GetFixtureContext(gomock.Any()).
		Return(eventcore.FixtureContext{}, nil).AnyTimes()
	competitors := NewMockCompetitorReader(ctrl)
	competitors.EXPECT().ListFixtureCompetitors(gomock.Any()).
		Return(nil, nil).AnyTimes()
	rndr := NewMockRenderer(ctrl)
	rndr.EXPECT().Render(gomock.Any()).Return("x").AnyTimes()

	return New(
		fixtures, repo, disp, context, competitors, rndr,
		mentions, channelID, testGuild, 0,
	)
}

// enrichedReaders builds the two enrichment readers, each expecting one call
// returning the given values.
func enrichedReaders(
	ctrl *gomock.Controller,
	context eventcore.FixtureContext,
	competitors []eventcore.FixtureCompetitor,
) (*MockFixtureContextReader, *MockCompetitorReader) {
	contextReader := NewMockFixtureContextReader(ctrl)
	contextReader.EXPECT().GetFixtureContext(gomock.Any()).Return(context, nil)
	competitorReader := NewMockCompetitorReader(ctrl)
	competitorReader.EXPECT().ListFixtureCompetitors(gomock.Any()).Return(competitors, nil)
	return contextReader, competitorReader
}

func TestNotifier_NotifyPending_SkipsAlreadySent(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("abc")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
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
	context := eventcore.FixtureContext{
		Competition: eventcore.Competition{Code: "ATH", Discipline: "Athletics"},
	}
	competitors := []eventcore.FixtureCompetitor{{
		Participant: eventcore.Participant{Name: "Usain Bolt", CountryISO: "JAM"},
		Result:      eventcore.Result{Outcome: eventcore.OutcomeMedalGold},
	}}

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
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

	contextReader, competitorReader := enrichedReaders(ctrl, context, competitors)
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
		fixtures, repo, disp, contextReader, competitorReader, rndr,
		mentions, defaultChan, testGuild, 0,
	)

	if err := n.NotifyPending(); err != nil {
		t.Fatalf("NotifyPending: %v", err)
	}
	if sentContent != "x\n<@u1> <@u2>" {
		t.Fatalf("content = %q, want body then mention footer", sentContent)
	}
	assertDispatchedAndSent(t, lastView, upserted)
}

func assertDispatchedAndSent(
	t *testing.T, lastView render.FixtureView, upserted []eventcore.Notification,
) {
	t.Helper()
	if lastView.Context.Competition.Code != "ATH" || len(lastView.Competitors) != 1 {
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

func TestNotifier_NotifyPending_NoMentions_NoSuffix(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("nopfx")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
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
		t.Fatalf("content = %q, want body without footer", sentContent)
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
	noPriorSent(repo)
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
		NewMockFixtureContextReader(ctrl),
		NewMockCompetitorReader(ctrl),
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

func TestNotifier_EditsInPlaceOnChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("new")

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetLatestSentNotificationByAlert(f.ID).
		Return(eventcore.Notification{
			AlertID:   f.ID,
			MessageID: "m1",
			Status:    eventcore.NotificationSent,
			Checksum:  "old",
		}, nil)
	var updated eventcore.Notification
	repo.EXPECT().
		UpsertNotification(gomock.Any()).
		Do(func(nt eventcore.Notification) { updated = nt }).
		Return(nil)

	disp := NewMockDispatcher(ctrl)
	var edited struct {
		channel, message string
	}
	disp.EXPECT().
		Edit(defaultChan, "m1", gomock.Any()).
		Do(func(channel, message, _ string) {
			edited.channel, edited.message = channel, message
		}).
		Return(nil)

	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().MentionsFor(testGuild, gomock.Any(), gomock.Any()).Return(nil, nil)

	n := newTestNotifier(
		ctrl, NewMockFixtureReader(ctrl), repo, disp, mentions, defaultChan,
	)

	n.On(f)
	if edited.channel != defaultChan || edited.message != "m1" {
		t.Fatalf("edit target mismatch: %+v", edited)
	}
	if updated.MessageID != "m1" || updated.Checksum != "new" ||
		updated.Status != eventcore.NotificationSent {
		t.Fatalf("edited notification not updated correctly: %+v", updated)
	}
}

func TestNotifier_NoOpWhenUnchanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("same")

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().
		GetLatestSentNotificationByAlert(f.ID).
		Return(eventcore.Notification{
			AlertID:   f.ID,
			MessageID: "m1",
			Status:    eventcore.NotificationSent,
			Checksum:  "same",
		}, nil)

	disp := NewMockDispatcher(ctrl) // neither Send nor Edit expected
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
		NewMockFixtureContextReader(ctrl), NewMockCompetitorReader(ctrl),
		NewMockRenderer(ctrl), NewMockMentionResolver(ctrl), "c", testGuild, 0,
	)
	if n.window == 0 {
		t.Fatal("default window must be non-zero")
	}
}
