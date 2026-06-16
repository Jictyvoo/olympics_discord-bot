package notifier

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Enrichment reads are best-effort: when results/competition/participants all
// fail, the notification still dispatches with a degraded view.
func TestNotifier_BuildView_EnrichmentErrors_StillDispatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("degraded")

	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().GetNotificationByChecksum("degraded").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	repo.EXPECT().UpsertNotification(gomock.Any()).Return(nil).Times(2)

	results := NewMockResultReader(ctrl)
	results.EXPECT().ListResultsByFixture(gomock.Any()).Return(nil, errors.New("results down"))
	comps := NewMockCompetitionReader(ctrl)
	// ErrNoRows is the non-warn branch; still degrades to a zero competition.
	comps.EXPECT().
		GetCompetitionByFixture(gomock.Any()).
		Return(eventcore.Competition{}, sql.ErrNoRows)
	parts := NewMockParticipantReader(ctrl)
	parts.EXPECT().ListParticipantsByFixture(gomock.Any()).Return(nil, errors.New("parts down"))

	rndr := NewMockRenderer(ctrl)
	rndr.EXPECT().Render(gomock.Any()).Return("body")
	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().MentionsFor(testGuild, gomock.Any(), gomock.Any()).Return(nil, nil)

	disp := NewMockDispatcher(ctrl)
	var sent bool
	disp.EXPECT().Send(defaultChan, "body").
		Do(func(_, _ string) { sent = true }).Return("m", nil)

	n := New(
		NewMockFixtureReader(ctrl), repo, disp, results, comps, parts, rndr,
		mentions, defaultChan, testGuild, 0,
	)
	n.On(f)

	if !sent {
		t.Fatal("expected dispatch despite enrichment errors")
	}
}

// A mention-resolver error must propagate before any pending row or dispatch.
func TestNotifier_MentionResolverError_Propagates(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("menterr")

	fixtures := NewMockFixtureReader(ctrl)
	fixtures.EXPECT().ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)
	repo := NewMockNotificationRepo(ctrl)
	repo.EXPECT().GetNotificationByChecksum("menterr").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	// No UpsertNotification, no Send: the error is raised while resolving mentions.

	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().MentionsFor(testGuild, gomock.Any(), gomock.Any()).
		Return(nil, errors.New("resolve fail"))

	n := newTestNotifier(ctrl, fixtures, repo, NewMockDispatcher(ctrl), mentions, defaultChan)

	if err := n.NotifyPending(); err == nil {
		t.Fatal("expected mention-resolver error to propagate")
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
