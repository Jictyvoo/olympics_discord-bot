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
