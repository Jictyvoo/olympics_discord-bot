package repolympicfetch

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/datasources/dsrest"
)

const apiURL = "https://sph-s-api.olympics.com/summer/schedules/api/ENG/schedule/day/"

type OlympicsFetcherImpl struct {
	ds dsrest.RESTDataSource
}

func NewOlympicsFetcher(ds dsrest.RESTDataSource) usecases.OlympicsFetcher {
	return OlympicsFetcherImpl{ds: ds}
}

func (repo OlympicsFetcherImpl) parseAPICompetitors(
	competitors []OlympicsAPIResponseCompetitors,
) (resultComp []entities.OlympicCompetitors) {
	resultComp = make([]entities.OlympicCompetitors, 0, len(competitors))
	for _, competitor := range competitors {
		newCompetitor := entities.OlympicCompetitors{
			Code:        competitor.Code,
			CountryCode: competitor.Noc,
			CountryInfo: entities.GetCountryByCode(competitor.Noc),
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
			EventName:      unit.EventUnitName,
			DisciplineName: unit.DisciplineName,
			Phase:          unit.PhaseName,
			StartAt:        unit.StartDate,
			EndAt:          unit.EndDate,
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

func (repo OlympicsFetcherImpl) FetchDataFromDay(day time.Time) ([]entities.OlympicEvent, error) {
	url := apiURL + day.Format(time.DateOnly)

	resp, err := repo.ds.Get(url)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonResp OlympicsAPIResponse
	err = decoder.Decode(&jsonResp)

	return repo.parseAPIResp(jsonResp), err
}
