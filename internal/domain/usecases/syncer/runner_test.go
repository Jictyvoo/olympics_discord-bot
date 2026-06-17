package syncer

import (
	"testing"
	"time"
)

func TestSyncWindow_SpansLookbackToLookahead(t *testing.T) {
	cases := []struct {
		name     string
		now      time.Time
		wantFrom time.Time
		wantTo   time.Time
	}{
		{
			name:     "midday spans yesterday to tomorrow",
			now:      time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC),
			wantFrom: time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC),
			wantTo:   time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC),
		},
		{
			// A 21:00 UTC-3 kickoff is 00:00 UTC the next day; the look-ahead must
			// reach that day so SyncRange fetches it before kickoff.
			name:     "late evening reaches next UTC day",
			now:      time.Date(2026, 6, 19, 20, 30, 0, 0, time.UTC),
			wantFrom: time.Date(2026, 6, 18, 20, 30, 0, 0, time.UTC),
			wantTo:   time.Date(2026, 6, 20, 20, 30, 0, 0, time.UTC),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			from, to := syncWindow(tc.now)
			if !from.Equal(tc.wantFrom) {
				t.Errorf("from = %v, want %v", from, tc.wantFrom)
			}
			if !to.Equal(tc.wantTo) {
				t.Errorf("to = %v, want %v", to, tc.wantTo)
			}
		})
	}
}

func TestNewRunner_DefaultsInterval(t *testing.T) {
	cases := []struct {
		name string
		in   time.Duration
		want time.Duration
	}{
		{"zero defaults", 0, defaultSyncIntervalMinutes * time.Minute},
		{"negative defaults", -time.Second, defaultSyncIntervalMinutes * time.Minute},
		{"positive kept", 30 * time.Second, 30 * time.Second},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRunner(nil, tc.in)
			if r.interval != tc.want {
				t.Errorf("interval = %v, want %v", r.interval, tc.want)
			}
		})
	}
}
