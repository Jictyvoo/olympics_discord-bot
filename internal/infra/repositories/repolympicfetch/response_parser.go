package repolympicfetch

import (
	"strings"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func (repo OlympicsFetcherImpl) parseAPICompetitors(
	competitors []OlympicsAPIResponseCompetitors,
) (resultComp []entities.OlympicCompetitors) {
	resultComp = make([]entities.OlympicCompetitors, 0, len(competitors))
	for _, competitor := range competitors {
		newCompetitor := entities.OlympicCompetitors{
			Code:        competitor.Code,
			CountryCode: competitor.Noc,
			Country:     entities.GetCountryByCode(competitor.Noc),
			Name:        competitor.Name,
		}

		resultComp = append(resultComp, newCompetitor)
	}

	return
}

func (repo OlympicsFetcherImpl) parseAPIResp(response OlympicsAPIResponse) []entities.OlympicEvent {
	events := make([]entities.OlympicEvent, 0, len(response.Units))
	for _, unit := range response.Units {
		newEvent := entities.OlympicEvent{
			EventName: unit.EventUnitName,
			Discipline: entities.Discipline{
				Name: unit.DisciplineName,
				Code: unit.DisciplineCode,
			},
			Phase:       unit.PhaseName,
			Gender:      entities.GenderOther,
			SessionCode: unit.UnitNum + "_#" + unit.SessionCode,
			StartAt:     unit.StartDate,
			EndAt:       unit.EndDate,
		}

		switch {
		case unit.GenderCode == "W":
			newEvent.Gender = entities.GenderFem
		case unit.GenderCode == "M":
			newEvent.Gender = entities.GenderMasc
		}
		switch strings.ToLower(unit.Status) {
		case string(entities.StatusScheduled):
			newEvent.Status = entities.StatusScheduled
		case string(entities.StatusFinished):
			newEvent.Status = entities.StatusFinished
		}

		newEvent.Competitors = repo.parseAPICompetitors(unit.Competitors)
		events = append(events, newEvent)
	}
	return events
}
