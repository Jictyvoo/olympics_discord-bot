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
		SaveEvent(event entities.OlympicEvent, competitorIDs []entities.Identifier) error
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

func (uc FetcherCacheUseCase) Run(date time.Time) (err error) {
	// Start fetching all elements
	var events []entities.OlympicEvent
	events, err = uc.fetcherRepo.FetchDataFromDay(date)
	if err != nil {
		return
	}

	// Start to insert on database
	for _, event := range events {
		competitorIDs := make([]entities.Identifier, 0, len(event.Competitors))
		for _, competitor := range event.Competitors {
			var compID entities.Identifier
			if compID, err = uc.storageRepo.SaveCompetitor(competitor); err != nil {
				return
			}
			competitorIDs = append(competitorIDs, compID)
		}

		if err = uc.storageRepo.SaveEvent(event, competitorIDs); err != nil {
			return
		}
	}

	return
}

func (uc FetcherCacheUseCase) FetchDisciplines() (err error) {
	// Start fetching all elements
	var disciplines []entities.Discipline
	if disciplines, err = uc.fetcherRepo.FetchDisciplines(); err != nil {
		return
	}

	// Start to insert on database
	err = uc.storageRepo.SaveDisciplines(disciplines)
	return
}
