package syncer

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
)

func TestSyncer_Persist_WritesAllCollectionsInFKSafeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	fixID := eventcore.NewID(eventcore.ProviderOlympics, "f1")
	delta := eventcore.SyncDelta{
		Competitions: []eventcore.Competition{
			{ID: eventcore.NewID(eventcore.ProviderOlympics, "c1")},
		},
		Seasons: []eventcore.Season{{ID: eventcore.NewID(eventcore.ProviderOlympics, "se1")}},
		Stages:  []eventcore.Stage{{ID: eventcore.NewID(eventcore.ProviderOlympics, "st1")}},
		Groups:  []eventcore.Group{{ID: eventcore.NewID(eventcore.ProviderOlympics, "g1")}},
		Venues:  []eventcore.Venue{{ID: eventcore.NewID(eventcore.ProviderOlympics, "v1")}},
		Participants: []eventcore.Participant{
			{ID: eventcore.NewID(eventcore.ProviderOlympics, "p1")},
		},
		Fixtures: []eventcore.Fixture{
			{ID: fixID, Participants: []eventcore.FixtureParticipant{{Role: "home"}}},
		},
		Results: []eventcore.Result{{FixtureID: fixID}},
		Standings: []eventcore.Standing{
			{StageID: eventcore.NewID(eventcore.ProviderOlympics, "st1")},
		},
	}

	// gomock.InOrder enforces the FK-safe write order, then Commit.
	gomock.InOrder(
		tx.EXPECT().UpsertCompetition(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertSeason(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertStage(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertGroup(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertVenue(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertParticipant(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertFixture(gomock.Any()).Return(nil),
		tx.EXPECT().
			UpsertFixtureParticipants(gomock.Any(), gomock.Any()).
			Return(nil),
		tx.EXPECT().UpsertResult(gomock.Any()).Return(nil),
		tx.EXPECT().UpsertStanding(gomock.Any()).Return(nil),
		tx.EXPECT().Commit().Return(nil),
	)
	tx.EXPECT().Rollback().Return(nil).AnyTimes()

	s := New(stubSet{}, repo, nil, t.Context())
	if err := s.persist(delta); err != nil {
		t.Fatalf("persist: %v", err)
	}
}

type recordingObserver struct{ got []eventcore.Fixture }

func (o *recordingObserver) On(f eventcore.Fixture) { o.got = append(o.got, f) }

func TestSyncer_Persist_EmitsFixturesAfterCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	fixID := eventcore.NewID(eventcore.ProviderOlympics, "f1")
	delta := eventcore.SyncDelta{
		Fixtures: []eventcore.Fixture{{ID: fixID}},
	}

	tx.EXPECT().UpsertFixture(gomock.Any()).Return(nil)
	tx.EXPECT().UpsertFixtureParticipants(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().Commit().Return(nil)
	tx.EXPECT().Rollback().Return(nil).AnyTimes()

	subject := services.NewSubject[eventcore.Fixture]()
	obs := &recordingObserver{}
	subject.Register(obs)

	s := New(stubSet{}, repo, subject, t.Context())
	if err := s.persist(delta); err != nil {
		t.Fatalf("persist: %v", err)
	}
	if len(obs.got) != 1 || obs.got[0].ID != fixID {
		t.Fatalf("expected one emitted fixture %v, got %v", fixID, obs.got)
	}
}

func TestSyncer_Persist_DoesNotEmitWhenCommitFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	delta := eventcore.SyncDelta{
		Fixtures: []eventcore.Fixture{{ID: eventcore.NewID(eventcore.ProviderOlympics, "f1")}},
	}

	tx.EXPECT().UpsertFixture(gomock.Any()).Return(nil)
	tx.EXPECT().UpsertFixtureParticipants(gomock.Any(), gomock.Any()).Return(nil)
	tx.EXPECT().Commit().Return(errors.New("commit failed"))
	tx.EXPECT().Rollback().Return(nil).AnyTimes()

	subject := services.NewSubject[eventcore.Fixture]()
	obs := &recordingObserver{}
	subject.Register(obs)

	s := New(stubSet{}, repo, subject, t.Context())
	if err := s.persist(delta); err == nil {
		t.Fatal("expected commit error")
	}
	if len(obs.got) != 0 {
		t.Fatalf("no fixtures should be emitted when commit fails, got %v", obs.got)
	}
}
