package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/mocks"
	"github.com/jictyvoo/olympics_data_fetcher/internal/entities"
)

func TestFetcherCacheUseCase_FetchDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const eventID = 42697
	targetDay := time.Date(2024, time.August, 4, 0, 0, 0, 0, time.UTC)
	expectedEvents := []entities.OlympicEvent{
		{
			ID:                  eventID,
			Competitors:         []entities.OlympicCompetitors{{Code: "comp1"}},
			ResultPerCompetitor: map[string]entities.Results{"comp1": {Mark: "100"}},
		},
	}

	mockFetcherRepo := mocks.NewMockOlympicsFetcher(ctrl)
	mockStorageRepo := mocks.NewMockAccessDatabaseRepository(ctrl)

	// Register expected elements on mock
	{
		mockFetcherRepo.EXPECT().FetchDataFromDay(gomock.Any()).Return(expectedEvents, nil)
		mockStorageRepo.EXPECT().SaveCompetitor(gomock.Any()).Return(entities.Identifier(5982), nil)
		mockStorageRepo.EXPECT().SaveEvent(gomock.Any(), gomock.Any()).Return(nil)
	}

	uc := NewFetcherCacheUseCase(mockFetcherRepo, mockStorageRepo)

	events, err := uc.FetchDay(targetDay)

	assert.Equal(t, expectedEvents, events)
	assert.Equal(t, nil, err)
}

func TestFetcherCacheUseCase_FetchDisciplines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFetcherRepo := mocks.NewMockOlympicsFetcher(ctrl)
	mockStorageRepo := mocks.NewMockAccessDatabaseRepository(ctrl)

	const disciplineID = 0x6070042
	expectedDisciplines := []entities.Discipline{
		{ID: disciplineID, Name: "Discipline 1"},
	}

	{
		mockFetcherRepo.EXPECT().FetchDisciplines().Return(expectedDisciplines, nil)
		mockStorageRepo.EXPECT().SaveDisciplines(gomock.Any()).Return(nil)
	}

	uc := NewFetcherCacheUseCase(mockFetcherRepo, mockStorageRepo)

	disciplines, err := uc.FetchDisciplines()

	assert.Equal(t, expectedDisciplines, disciplines)
	assert.Equal(t, nil, err)
}
