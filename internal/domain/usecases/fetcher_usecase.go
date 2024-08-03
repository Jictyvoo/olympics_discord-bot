package usecases

import (
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

type (
	OlympicsFetcher interface {
		FetchDataFromDay(day time.Time) ([]entities.OlympicEvent, error)
		FetchDisciplines() ([]entities.Discipline, error)
	}
	AccessDatabaseRepository interface {
		InsertCountries(countries []entities.CountryInfo) ([]entities.Identifier, error)
		InsertCountry(country entities.CountryInfo) (entities.Identifier, error)
		SaveCompetitor(competitors entities.OlympicCompetitors) (entities.Identifier, error)
		SaveDisciplines(disciplineList []entities.Discipline) error
		SaveEvent(
			event entities.OlympicEvent,
			competitorResultsByIDs map[entities.Identifier]*entities.Results,
		) error
	}
)

type FetcherCacheUseCase struct {
	fetcherRepo OlympicsFetcher
	storageRepo AccessDatabaseRepository
}

func NewFetcherCacheUseCase(
	fetcherRepo OlympicsFetcher,
	storageRepo AccessDatabaseRepository,
) FetcherCacheUseCase {
	return FetcherCacheUseCase{fetcherRepo: fetcherRepo, storageRepo: storageRepo}
}

func (uc FetcherCacheUseCase) FetchDay(date time.Time) (events []entities.OlympicEvent, err error) {
	// Start fetching all elements
	events, err = uc.fetcherRepo.FetchDataFromDay(date)
	if err != nil {
		return
	}

	// Start to insert on database
	for _, event := range events {
		resultPerCompetitorID := make(
			map[entities.Identifier]*entities.Results, len(event.Competitors),
		)
		for _, competitor := range event.Competitors {
			var compID entities.Identifier
			if compID, err = uc.storageRepo.SaveCompetitor(competitor); err != nil {
				return
			}

			var compResult *entities.Results
			if fetchedResult, ok := event.ResultPerCompetitor[competitor.Code]; ok {
				compResult = &fetchedResult
			}
			resultPerCompetitorID[compID] = compResult
		}

		if err = uc.storageRepo.SaveEvent(event, resultPerCompetitorID); err != nil {
			return
		}
	}

	return
}

func (uc FetcherCacheUseCase) FetchDisciplines() (disciplines []entities.Discipline, err error) {
	// Start fetching all elements
	if disciplines, err = uc.fetcherRepo.FetchDisciplines(); err != nil {
		return
	}

	// Start to insert on database
	err = uc.storageRepo.SaveDisciplines(disciplines)
	return
}
