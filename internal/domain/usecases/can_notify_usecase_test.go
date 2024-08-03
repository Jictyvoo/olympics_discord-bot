package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func TestCanNotifyUseCase_timeDiffAllowed(t *testing.T) {
	fixedNow := time.Now()

	tests := []struct {
		name            string
		event           entities.OlympicEvent
		allowedTimeDiff time.Duration
		expected        bool
	}{
		{
			name: "Event within allowed time diff",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-5 * time.Hour),
				EndAt:   fixedNow.Add(5 * time.Hour),
			},
			allowedTimeDiff: 10 * time.Hour,
			expected:        true,
		},
		{
			name: "Event exactly on the boundary",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-5 * time.Hour),
				EndAt:   fixedNow.Add(5 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        true,
		},
		{
			name: "Event outside allowed time diff",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-6 * time.Hour),
				EndAt:   fixedNow.Add(6 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
		{
			name: "Event in the past",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(-12 * time.Hour),
				EndAt:   fixedNow.Add(-6 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
		{
			name: "Event in the future",
			event: entities.OlympicEvent{
				StartAt: fixedNow.Add(6 * time.Hour),
				EndAt:   fixedNow.Add(12 * time.Hour),
			},
			allowedTimeDiff: 5 * time.Hour,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				uc := CanNotifyUseCase{
					allowedTimeDiff: tt.allowedTimeDiff, timeNow: func() time.Time {
						return fixedNow
					},
				}
				result := uc.timeDiffAllowed(tt.event)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
