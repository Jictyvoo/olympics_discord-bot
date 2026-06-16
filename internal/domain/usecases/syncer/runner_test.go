package syncer

import (
	"testing"
	"time"
)

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
