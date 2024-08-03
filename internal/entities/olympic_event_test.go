package entities

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestSHAIdentifier(t *testing.T) {
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
