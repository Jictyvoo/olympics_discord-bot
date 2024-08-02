package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func TestOlympicEventManager_GenContent(t *testing.T) {
	manager, err := NewOlympicEventManager([]string{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	event := entities.OlympicEvent{
		Discipline: entities.Discipline{
			Name: "Athletics",
			Code: "ATH",
		},
		EventName: "100m Final",
		Status:    "Completed",
		Phase:     "Final",
		Gender:    entities.GenderMasc,
		StartAt:   time.Date(2024, 7, 30, 18, 0, 0, 0, time.UTC),
		EndAt:     time.Date(2024, 7, 30, 18, 10, 0, 0, time.UTC),
		Competitors: []entities.OlympicCompetitors{
			{Country: entities.CountryInfo{ISOCode: [2]string{"us"}}, Name: "John Doe", Code: "25"},
			{
				Country: entities.CountryInfo{ISOCode: [2]string{"ca"}},
				Name:    "Alex Smith",
				Code:    "28",
			},
		},
		ResultPerCompetitor: map[string]entities.Results{
			"25": {MedalType: entities.MedalLoser, Mark: "52.5"},
			"28": {MedalType: entities.MedalWinner},
		},
	}

	expectedOutput := `
# :athletic_shoe: Athletics
**Event:** 100m Final - Completed
**Phase:** Final
**Gender:** Male
**Start:** <t:1722362400:R>
**End:** <t:1722363000:R>
**Competitors:**
- :flag_us: John Doe #52.5 (loser)
- :flag_ca: Alex Smith (winner)
`

	content := manager.genContent(event)
	assert.Equal(t, expectedOutput, content)
}

func TestOlympicEventManager_NormalizeEvent4Notification(t *testing.T) {
	tests := []struct {
		name                string
		watchCountries      []string
		event               *entities.OlympicEvent
		expectedResult      bool
		expectedCompetitors []entities.OlympicCompetitors
	}{
		{
			name:           "No watched countries",
			watchCountries: []string{},
			event: &entities.OlympicEvent{
				Competitors: []entities.OlympicCompetitors{
					{Code: "01", Country: entities.CountryInfo{IOCCode: "USA"}},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"USA": {MedalType: entities.MedalGold},
				},
			},
			expectedResult: true,
		},
		{
			name:           "No matching competitors",
			watchCountries: []string{"BRA"},
			event: &entities.OlympicEvent{
				Competitors: []entities.OlympicCompetitors{
					{Code: "01", Country: entities.CountryInfo{IOCCode: "USA"}},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"USA": {MedalType: entities.MedalGold},
				},
			},
			expectedResult: false,
		},
		{
			name:           "Some matching competitors",
			watchCountries: []string{"USA"},
			event: &entities.OlympicEvent{
				Competitors: []entities.OlympicCompetitors{
					{Code: "01", Country: entities.CountryInfo{IOCCode: "USA"}},
					{Code: "02", Country: entities.CountryInfo{IOCCode: "BRA"}},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"USA": {MedalType: entities.MedalGold},
				},
			},
			expectedResult: true,
		},
		{
			name:           "Competitors with results",
			watchCountries: []string{"BRA"},
			event: &entities.OlympicEvent{
				Competitors: []entities.OlympicCompetitors{
					{Code: "01", Country: entities.CountryInfo{IOCCode: "USA"}},
					{Code: "02", Country: entities.CountryInfo{IOCCode: "BRA"}},
				},
				ResultPerCompetitor: map[string]entities.Results{
					"01": {MedalType: entities.MedalGold},
				},
			},
			expectedResult: true,
		},
		{
			name:           "Competitors list more than 4",
			watchCountries: []string{"BRA"},
			event: &entities.OlympicEvent{
				Competitors: []entities.OlympicCompetitors{
					{Code: "01", Country: entities.CountryInfo{IOCCode: "BRA"}},
					{Code: "02", Country: entities.CountryInfo{IOCCode: "BRA"}},
					{Code: "03", Country: entities.CountryInfo{IOCCode: "BRA"}},
					{Code: "04", Country: entities.CountryInfo{IOCCode: "BRA"}},
					{Code: "05", Country: entities.CountryInfo{IOCCode: "BRA"}},
				},
				ResultPerCompetitor: map[string]entities.Results{},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				oen := OlympicEventManager{watchCountries: tt.watchCountries}
				result := oen.NormalizeEvent4Notification(tt.event)
				assert.Equal(t, tt.expectedResult, result)
				// assert.Equal(t, tt.event.Competitors, tt.expectedCompetitors)
			},
		)
	}
}
