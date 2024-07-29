package entities

import (
	"testing"
)

func TestSHAIdentifier(t *testing.T) {
	tests := []struct {
		name  string
		event OlympicEvent
	}{
		{
			name: "Test Case 1",
			event: OlympicEvent{
				DisciplineName: "Athletics",
				Gender:         1,
				Phase:          "Final",
				EventName:      "100m",
			},
		},
		{
			name: "Test Case 2",
			event: OlympicEvent{
				DisciplineName: "Swimming",
				Gender:         0,
				Phase:          "Semifinal",
				EventName:      "200m Freestyle",
			},
		},
		{
			name: "Test Case 3",
			event: OlympicEvent{
				DisciplineName: "Basketball",
				Gender:         1,
				Phase:          "Quarterfinal",
				EventName:      "Men's Basketball",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got1 := tt.event.SHAIdentifier()
				got2 := tt.event.SHAIdentifier()
				if got1 != got2 {
					t.Errorf(
						"Expected SHAIdentifier to produce the same SHAIdentifier. Got: %v, Want: %v",
						got1,
						got2,
					)
				}
			},
		)
	}
}
