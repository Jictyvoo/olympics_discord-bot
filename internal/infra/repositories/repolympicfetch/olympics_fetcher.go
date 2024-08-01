package repolympicfetch

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/datasources/dsrest"
)

type OlympicsFetcherImpl struct {
	ds   dsrest.RESTDataSource
	lang entities.Language
}

func NewOlympicsFetcher(
	apiLocale entities.Language,
	ds dsrest.RESTDataSource,
) usecases.OlympicsFetcher {
	return OlympicsFetcherImpl{lang: apiLocale, ds: ds}
}

func (repo OlympicsFetcherImpl) FetchDataFromDay(day time.Time) ([]entities.OlympicEvent, error) {
	const apiURL = "https://sph-s-api.olympics.com/summer/schedules/api/%s/schedule/day/%s"
	url := fmt.Sprintf(apiURL, repo.lang.Code, day.Format(time.DateOnly))

	resp, err := repo.ds.Get(url)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonResp OlympicsAPIResponse
	err = decoder.Decode(&jsonResp)

	return repo.parseAPIResp(jsonResp), err
}

func (repo OlympicsFetcherImpl) FetchWatchOn() ([]string, error) {
	const apiURL = "https://sph-i-api.olympics.com/summer/info/api/%s/mrh"
	resp, err := repo.ds.Get(fmt.Sprintf(apiURL, repo.lang.Code))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonResp WatchOnResp
	err = decoder.Decode(&jsonResp)

	urls := make([]string, 0, len(jsonResp.MrhItems))
	for _, item := range jsonResp.MrhItems {
		if item.Url != "" {
			urls = append(urls, item.Url)
		}
	}

	return urls, nil
}

func (repo OlympicsFetcherImpl) FetchDisciplines() ([]entities.Discipline, error) {
	const apiURL = "https://sph-i-api.olympics.com/summer/info/api/%s/disciplines"
	resp, err := repo.ds.Get(fmt.Sprintf(apiURL, repo.lang.Code))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonResp []DisciplineResp
	err = decoder.Decode(&jsonResp)

	resultDisciplines := make([]entities.Discipline, 0, len(jsonResp))
	for _, item := range jsonResp {
		discipline := entities.Discipline{
			Code:         item.Code,
			Name:         item.Description,
			Description:  item.Description,
			IsSport:      item.IsSport,
			IsParalympic: item.IsParalympic,
		}
		resultDisciplines = append(resultDisciplines, discipline)
	}

	return resultDisciplines, nil
}

func (repo OlympicsFetcherImpl) FetchCompetitionDays() ([]time.Time, error) {
	const apiURL = "https://sph-s-api.olympics.com/summer/schedules/api/%s/competitiondays"
	resp, err := repo.ds.Get(fmt.Sprintf(apiURL, repo.lang.Code))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var jsonResp struct {
		Days []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"days"`
	}
	err = decoder.Decode(&jsonResp)

	days := make([]time.Time, 0, len(jsonResp.Days))
	for _, item := range jsonResp.Days {
		if parsed, parsErr := time.Parse(time.DateOnly, item.Id); parsErr == nil {
			days = append(days, parsed)
		}
	}

	return days, nil
}
