package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		expected  time.Time
		wantErr   bool
	}{
		{
			name:      "RFC3339 with time zone",
			timestamp: "2024-08-03 14:22:00 -0500 CDT",
			expected: time.Date(
				2024, time.August, 3, 14, 22, 0, 0,
				time.FixedZone("CDT", -5*3600),
			),
		},
		{
			name:      "RFC3339 with offset and time zone",
			timestamp: "2024-08-06 15:35:00 +0200 +0200",
			expected:  time.Date(2024, time.August, 6, 15, 35, 0, 0, time.FixedZone("", 2*3600)),
		},
		{
			name:      "Datetime with offset",
			timestamp: "2024-08-03 14:22:00 -0500 -0500",
			expected:  time.Date(2024, time.August, 3, 14, 22, 0, 0, time.FixedZone("", -5*3600)),
		},
		{
			name:      "RFC3339 without offset",
			timestamp: "2024-08-03T14:22:00Z",
			expected:  time.Date(2024, time.August, 3, 14, 22, 0, 0, time.UTC),
		},
		{
			name:      "Datetime with time zone",
			timestamp: "2024-08-06 16:00:00 +0000 UTC",
			expected:  time.Date(2024, time.August, 6, 16, 0, 0, 0, time.UTC),
		},
		{
			name:      "Invalid format",
			timestamp: "invalid-timestamp",
			expected:  time.Time{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result, err := ParseTimestamp(tt.timestamp)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			},
		)
	}
}

func TestEnsureTime(t *testing.T) {
	dynamicNow := time.Now()
	fixedNow := time.Date(2024, 8, 3, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    *time.Time
		duration time.Duration
		expected time.Time
	}{
		{
			name:     "Non-nil and non-zero value",
			value:    &fixedNow,
			duration: 5 * time.Hour,
			expected: fixedNow,
		},
		{
			name:     "Nil value",
			value:    nil,
			duration: 5 * time.Hour,
			expected: time.Date(
				dynamicNow.Year(), dynamicNow.Month(), dynamicNow.Day(),
				0, 0, 0, 0, dynamicNow.Location(),
			).Add(5 * time.Hour),
		},
		{
			name:     "Zero value",
			value:    &time.Time{},
			duration: 5 * time.Hour,
			expected: time.Date(
				dynamicNow.Year(), dynamicNow.Month(), dynamicNow.Day(),
				0, 0, 0, 0, dynamicNow.Location(),
			).Add(5 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := EnsureTime(tt.value, tt.duration)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
