package syncer

import (
	"context"
	"errors"
	"iter"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
)

// stubStrategy / stubSet implement provider.Strategy / provider.Set, which are
// declared in the provider package (not in this package's interfaces.go), so
// they are kept as hand-written stubs rather than generated mocks.
type stubStrategy struct {
	code  eventcore.ProviderID
	delta eventcore.SyncDelta
	err   error
}

func (s stubStrategy) Code() eventcore.ProviderID { return s.code }
func (stubStrategy) DisplayName() string          { return "stub" }
func (s stubStrategy) SyncFixturesByDate(context.Context, time.Time) (eventcore.SyncDelta, error) {
	return s.delta, s.err
}

func (stubStrategy) SyncFixtureResults(
	context.Context, eventcore.Fixture,
) (eventcore.SyncDelta, error) {
	return eventcore.SyncDelta{}, nil
}

type stubSet struct{ strategies []provider.Strategy }

//nolint:ireturn // factory returning consumer interface by design
func (s stubSet) Get(eventcore.ProviderID) (provider.Strategy, error) {
	return nil, errors.New("not used")
}

func (s stubSet) Enabled() iter.Seq[provider.Strategy] {
	return func(yield func(provider.Strategy) bool) {
		for _, st := range s.strategies {
			if !yield(st) {
				return
			}
		}
	}
}

// repoWithTx wires a MockRepository so that Begin returns the given MockTx.
func repoWithTx(ctrl *gomock.Controller) (*MockRepository, *MockTx) {
	repo := NewMockRepository(ctrl)
	tx := NewMockTx(ctrl)
	repo.EXPECT().Begin().Return(tx, nil).AnyTimes()
	return repo, tx
}

func TestSyncer_SyncDay_PersistsAndCheckpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	fixID := eventcore.NewID(eventcore.ProviderOlympics, "f1")
	strategy := stubStrategy{
		code: eventcore.ProviderOlympics,
		delta: eventcore.SyncDelta{
			Cursor: "2024-08-04",
			Participants: []eventcore.Participant{
				{ID: eventcore.NewID(eventcore.ProviderOlympics, "p1")},
			},
			Fixtures: []eventcore.Fixture{
				{
					ID:           fixID,
					Name:         "100m final",
					Participants: []eventcore.FixtureParticipant{{Role: "athlete"}},
				},
			},
			Results: []eventcore.Result{{FixtureID: fixID}},
		},
	}

	tx.EXPECT().UpsertParticipant(gomock.Any()).Return(nil).Times(1)
	tx.EXPECT().UpsertFixture(gomock.Any()).Return(nil).Times(1)
	tx.EXPECT().
		UpsertFixtureParticipants(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)
	tx.EXPECT().UpsertResult(gomock.Any()).Return(nil).Times(1)
	tx.EXPECT().Commit().Return(nil).Times(1)
	tx.EXPECT().Rollback().Return(nil).AnyTimes()
	repo.EXPECT().
		SaveCursor(eventcore.ProviderOlympics, gomock.Any(), "2024-08-04").
		Return(nil).
		Times(1)
	// RecordError must NOT be called on the happy path: no EXPECT() => any call fails.

	s := New(stubSet{strategies: []provider.Strategy{strategy}}, repo, nil, t.Context())

	if err := s.SyncDay(time.Date(2024, 8, 4, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("SyncDay: %v", err)
	}
}

func TestSyncer_SyncDay_FetchErrorRecordsAndContinues(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	good := stubStrategy{code: eventcore.ProviderFIFA, delta: eventcore.SyncDelta{Cursor: "ok"}}
	bad := stubStrategy{code: eventcore.ProviderOlympics, err: errors.New("upstream 500")}

	tx.EXPECT().Commit().Return(nil).AnyTimes()
	tx.EXPECT().Rollback().Return(nil).AnyTimes()
	repo.EXPECT().
		RecordError(eventcore.ProviderOlympics, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)
	repo.EXPECT().
		SaveCursor(eventcore.ProviderFIFA, gomock.Any(), "ok").
		Return(nil).
		Times(1)

	s := New(stubSet{strategies: []provider.Strategy{bad, good}}, repo, nil, t.Context())

	if err := s.SyncDay(time.Now()); err == nil {
		t.Fatal("expected joined error")
	}
}

func TestSyncer_SyncDay_PersistErrorRecordsAndReturns(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	strategy := stubStrategy{
		code: eventcore.ProviderOlympics,
		delta: eventcore.SyncDelta{
			Cursor:   "c",
			Fixtures: []eventcore.Fixture{{ID: eventcore.NewID(eventcore.ProviderOlympics, "f1")}},
		},
	}

	tx.EXPECT().
		UpsertFixture(gomock.Any()).
		Return(errors.New("constraint violation"))
	tx.EXPECT().Rollback().Return(nil).Times(1)
	tx.EXPECT().Commit().Return(nil).MaxTimes(0)
	repo.EXPECT().
		RecordError(eventcore.ProviderOlympics, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)
	// SaveCursor must NOT be called when persist fails: no EXPECT() => any call fails.

	s := New(stubSet{strategies: []provider.Strategy{strategy}}, repo, nil, t.Context())

	if err := s.SyncDay(time.Now()); err == nil {
		t.Fatal("expected error from persist failure")
	}
}

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

// recordingObserver captures fixtures emitted to it.
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
