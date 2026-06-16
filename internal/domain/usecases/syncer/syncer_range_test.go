package syncer

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
)

func TestSyncer_SyncRange_PersistsEachDay(t *testing.T) {
	cases := []struct {
		name     string
		from, to time.Time
		wantDays int
	}{
		{
			name:     "single day",
			from:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			wantDays: 1,
		},
		{
			name:     "three days",
			from:     time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC),
			wantDays: 3,
		},
		{
			name:     "from after to yields no days",
			from:     time.Date(2026, 6, 13, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			wantDays: 0,
		},
		{
			name:     "sub-day times still count one day",
			from:     time.Date(2026, 6, 11, 23, 30, 0, 0, time.UTC),
			to:       time.Date(2026, 6, 11, 1, 0, 0, 0, time.UTC),
			wantDays: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo, tx := repoWithTx(ctrl)
			tx.EXPECT().Commit().Return(nil).Times(tc.wantDays)
			tx.EXPECT().Rollback().Return(nil).AnyTimes()
			repo.EXPECT().
				SaveCursor(eventcore.ProviderFIFA, gomock.Any(), gomock.Any()).
				Return(nil).
				Times(tc.wantDays)

			strategy := stubStrategy{code: eventcore.ProviderFIFA}
			s := New(stubSet{strategies: []provider.Strategy{strategy}}, repo, nil, t.Context())

			if err := s.SyncRange(tc.from, tc.to); err != nil {
				t.Fatalf("SyncRange: %v", err)
			}
		})
	}
}

func TestSyncer_SyncRange_JoinsDayErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo, tx := repoWithTx(ctrl)
	tx.EXPECT().Rollback().Return(nil).AnyTimes()
	tx.EXPECT().Commit().Return(nil).AnyTimes()
	// Every day's provider fetch fails, so each records an error.
	repo.EXPECT().
		RecordError(eventcore.ProviderFIFA, gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)

	bad := stubStrategy{code: eventcore.ProviderFIFA, err: errors.New("upstream 500")}
	s := New(stubSet{strategies: []provider.Strategy{bad}}, repo, nil, t.Context())

	from := time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	if err := s.SyncRange(from, to); err == nil {
		t.Fatal("expected joined error across failing days")
	}
}
