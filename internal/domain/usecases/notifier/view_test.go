package notifier

import (
	"database/sql"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Enrichment reads are best-effort: when the context and competitor reads both
// fail, the notification still dispatches with a degraded view.
func TestNotifier_BuildView_EnrichmentErrors_StillDispatches(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("degraded")

	repo := NewMockNotificationRepo(ctrl)
	noPriorSent(repo)
	repo.EXPECT().GetNotificationByChecksum("degraded").
		Return(eventcore.Notification{}, sql.ErrNoRows)
	repo.EXPECT().UpsertNotification(gomock.Any()).Return(nil).Times(2)

	context := NewMockFixtureContextReader(ctrl)
	// ErrNoRows is the non-warn branch; still degrades to a zero context.
	context.EXPECT().
		GetFixtureContext(gomock.Any()).
		Return(eventcore.FixtureContext{}, sql.ErrNoRows)
	competitors := NewMockCompetitorReader(ctrl)
	competitors.EXPECT().
		ListFixtureCompetitors(gomock.Any()).
		Return(nil, errors.New("competitors down"))

	rndr := NewMockRenderer(ctrl)
	rndr.EXPECT().Render(gomock.Any()).Return("body")
	mentions := NewMockMentionResolver(ctrl)
	mentions.EXPECT().MentionsFor(testGuild, gomock.Any(), gomock.Any()).Return(nil, nil)

	disp := NewMockDispatcher(ctrl)
	var sent bool
	disp.EXPECT().Send(defaultChan, "body").
		Do(func(_, _ string) { sent = true }).Return("m", nil)

	n := New(
		NewMockFixtureReader(ctrl), repo, disp, context, competitors, rndr,
		mentions, channelRouter{fallback: defaultChan}, testGuild, 0,
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
	noPriorSent(repo)
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
