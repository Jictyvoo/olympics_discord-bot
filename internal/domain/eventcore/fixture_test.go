package entities

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOlympicEvent_SHAIdentifier(t *testing.T) {
	testCompetitors := mockOlympicCompetitors()
	tests := [...]struct {
		name  string
		event OlympicEvent
	}{
		{
			name: "Simple Event without Competitors",
			event: OlympicEvent{
				ID:                  58964455,
				EventName:           "200m Freestyle",
				Discipline:          Discipline{Name: "Swimming"},
				Phase:               "Semifinal",
				Gender:              GenderMasc,
				SessionCode:         "#56__EF4",
				StartAt:             time.Now(),
				EndAt:               time.Now().Add(time.Hour),
				Status:              StatusScheduled,
				Competitors:         nil,
				ResultPerCompetitor: nil,
			},
		},
		{
			name: "With 2 competitors",
			event: OlympicEvent{
				Discipline: Discipline{Name: "Athletics"},
				Gender:     1,
				Phase:      "Final",
				EventName:  "100m",
				Competitors: []OlympicCompetitors{
					testCompetitors[0],
					testCompetitors[1],
				},
			},
		},
		{
			name: "With result per competitor",
			event: OlympicEvent{
				ID:          25564329,
				EventName:   "Men's Basketball",
				Discipline:  Discipline{Name: "Basketball"},
				Phase:       "Quarterfinal",
				Gender:      GenderFem,
				SessionCode: "#54__4EDF",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(time.Hour),
				Status:      StatusFinished,
				Competitors: testCompetitors[:],
				ResultPerCompetitor: func() map[string]Results {
					resultList := make(map[string]Results, len(testCompetitors))
					medalList := []medalType{
						MedalNoMedal, MedalGold, MedalSilver,
						MedalBronze, MedalWinner, MedalLoser,
					}
					for index, competitor := range testCompetitors {
						resultList[competitor.Code] = Results{
							Position:  strconv.Itoa(index),
							Mark:      strconv.Itoa(rand.Int() % 110),
							MedalType: medalList[rand.Intn(len(medalList))],
							Irm:       "",
						}
					}
					return resultList
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				const totalChecks = 11
				originalResult := tt.event.SHAIdentifier()
				for range totalChecks {
					time.Sleep(time.Millisecond)
					if newResult := tt.event.SHAIdentifier(); newResult != originalResult {
						t.Fatalf(
							"Expected SHAIdentifier to produce the same SHAIdentifier. Got: %v, Want: %v",
							originalResult,
							newResult,
						)
					}
				}
			},
		)
	}
}

func TestOlympicEvent_Normalize(t *testing.T) {
	mockDates := [...]time.Time{
		time.Date(2024, time.August, 3, 12, 0, 0, 0, time.Local),
		time.Date(2024, time.August, 3, 15, 0, 0, 0, time.Local),
		time.Date(2024, time.August, 3, 15, 0, 0, 0, time.UTC),
		time.Date(2024, time.August, 3, 18, 0, 0, 0, time.UTC),
	}
	mockCompetitors := mockOlympicCompetitors()

	tests := []struct {
		name                string
		event               OlympicEvent
		expectedStartAt     time.Time
		expectedEndAt       time.Time
		expectedCompetitors []OlympicCompetitors
	}{
		{
			name: "Basic normalization and sorting",
			event: OlympicEvent{
				StartAt: mockDates[0],
				EndAt:   mockDates[1],
				Competitors: []OlympicCompetitors{
					mockCompetitors[1], mockCompetitors[0], mockCompetitors[2], mockCompetitors[3],
				},
				ResultPerCompetitor: map[string]Results{
					mockCompetitors[0].Code: {MedalType: MedalGold, Mark: "10.0"},
					mockCompetitors[1].Code: {MedalType: MedalSilver, Mark: "9.0"},
					mockCompetitors[2].Code: {MedalType: MedalBronze, Mark: "8.0"},
				},
			},
			expectedStartAt: mockDates[2],
			expectedEndAt:   mockDates[3],
			expectedCompetitors: []OlympicCompetitors{
				mockCompetitors[0], mockCompetitors[1], mockCompetitors[2], mockCompetitors[3],
			},
		},
		{
			name: "Sorting with no medals",
			event: OlympicEvent{
				StartAt: mockDates[0],
				EndAt:   mockDates[1],
				Competitors: []OlympicCompetitors{
					mockCompetitors[0], mockCompetitors[1], mockCompetitors[2],
				},
				ResultPerCompetitor: map[string]Results{
					mockCompetitors[0].Code: {MedalType: MedalNoMedal, Mark: "10"},
					mockCompetitors[1].Code: {MedalType: MedalNoMedal, Mark: "9"},
					mockCompetitors[2].Code: {MedalType: MedalNoMedal, Mark: "8.0"},
				},
			},
			expectedStartAt: mockDates[2],
			expectedEndAt:   mockDates[3],
			expectedCompetitors: []OlympicCompetitors{
				mockCompetitors[0], mockCompetitors[1], mockCompetitors[2],
			},
		},
		{
			name: "Sorting based on a time mark without a competitor result",
			event: OlympicEvent{
				StartAt: mockDates[0],
				EndAt:   mockDates[1],
				Competitors: []OlympicCompetitors{
					mockCompetitors[5], mockCompetitors[6], mockCompetitors[4],
				},
				ResultPerCompetitor: map[string]Results{
					mockCompetitors[5].Code: {MedalType: MedalNoMedal, Mark: "9:16.28"},
					mockCompetitors[6].Code: {MedalType: MedalNoMedal, Mark: "9:01.78"},
				},
			},
			expectedStartAt: mockDates[2],
			expectedEndAt:   mockDates[3],
			expectedCompetitors: []OlympicCompetitors{
				mockCompetitors[5], mockCompetitors[6], mockCompetitors[4],
			},
		},
		{
			name: "Sorting with mixed medals and marks",
			event: OlympicEvent{
				StartAt: mockDates[0],
				EndAt:   mockDates[1],
				Competitors: []OlympicCompetitors{
					mockCompetitors[0], mockCompetitors[1], mockCompetitors[2],
				},
				ResultPerCompetitor: map[string]Results{
					mockCompetitors[0].Code: {MedalType: MedalNoMedal, Mark: "10.0"},
					mockCompetitors[1].Code: {MedalType: MedalBronze, Mark: "8.0"},
					mockCompetitors[2].Code: {MedalType: MedalSilver, Mark: "9.0"},
				},
			},
			expectedStartAt: mockDates[2],
			expectedEndAt:   mockDates[3],
			expectedCompetitors: []OlympicCompetitors{
				mockCompetitors[2], mockCompetitors[1], mockCompetitors[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.event.Normalize()

				assert.Equal(t, tt.expectedStartAt, tt.event.StartAt)
				assert.Equal(t, tt.expectedEndAt, tt.event.EndAt)
				assert.Equal(t, tt.expectedCompetitors, tt.event.Competitors)
			},
		)
	}
}
