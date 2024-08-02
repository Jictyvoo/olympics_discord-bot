package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func TestGenContent(t *testing.T) {
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
			"25": {MedalType: entities.MedalLoser},
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
- :flag_us: John Doe (loser)
- :flag_ca: Alex Smith (winner)
`

	content := manager.genContent(event)
	assert.Equal(t, expectedOutput, content)
}
