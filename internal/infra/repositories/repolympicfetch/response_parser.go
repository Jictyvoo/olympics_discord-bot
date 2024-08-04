package repolympicfetch

import (
	"strings"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func (repo OlympicsFetcherImpl) parseAPIResults(
	competitors []OlympicsAPIResponseCompetitors,
) (resultComp map[string]entities.Results) {
	resultComp = make(map[string]entities.Results, len(competitors))
	for _, competitor := range competitors {
		competitorResult := entities.Results{
			Position: competitor.Results.Position,
			Mark:     competitor.Results.Mark,
			Irm:      competitor.Results.Irm,
		}

		conciseResult := competitor.Results.WinnerLoserTie
		if competitor.Results.MedalType != "" {
			conciseResult = competitor.Results.MedalType
		}
		competitorResult.MedalType = entities.Medal(conciseResult)

		if competitorResult.MedalType != entities.MedalNoMedal || competitorResult.Mark != "" {
			resultComp[competitor.Code] = competitorResult
		}
	}

	return
}

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
	groupMap := make(map[string]OlympicsAPIResponseGroup, len(response.Groups))
	for _, group := range response.Groups {
		groupMap[group.Id] = group
	}
	events := make([]entities.OlympicEvent, 0, len(response.Units))
	for _, unit := range response.Units {
		unitGroup := groupMap[unit.GroupId]
		newEvent := entities.OlympicEvent{
			EventName: unit.EventUnitName,
			Discipline: entities.Discipline{
				Name: unit.DisciplineName,
				Code: unit.DisciplineCode,
			},
			Phase:       unit.PhaseName,
			Gender:      entities.GenderOther,
			SessionCode: unit.UnitNum + "_#" + unit.SessionCode,
			HasMedal:    unitGroup.HasMedals,
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
		if newEvent.Status == "" && unitGroup.IsLive {
			newEvent.Status = entities.StatusOngoing
		}

		newEvent.Competitors = repo.parseAPICompetitors(unit.Competitors)
		newEvent.ResultPerCompetitor = repo.parseAPIResults(unit.Competitors)
		events = append(events, newEvent)
	}
	return events
}
