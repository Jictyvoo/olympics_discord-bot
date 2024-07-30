package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/bootstrap"
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
)

func main() {
	conf := bootstrap.Config()
	db := bootstrap.OpenDatabase()
	defer db.Close()

	inj := remy.NewInjector(remy.Config{})
	remy.RegisterInstance(inj, db)
	bootstrap.DoInjections(inj, conf)

	/*repo := remy.Get[usecases.AccessDatabaseRepository](inj)
	if _, err := repo.InsertCountries(entities.GetCountryList()); err != nil {
		return
	}*/

	startDate := time.Date(2024, time.July, 24, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, time.August, 12, 0, 0, 0, 0, time.UTC)
	for date := startDate; date.Before(endDate); date = date.Add(24 * time.Hour) {
		fmt.Printf("Fetching %s\n", date.Format("2006-01-02"))
		uc, err := remy.DoGet[usecases.FetcherCacheUseCase](inj)
		if err != nil {
			slog.Error("Failed to fetch Olympics", slog.String("error", err.Error()))
			continue
		}

		if err = uc.Run(date); err != nil {
			slog.Error("Error fetching data from day", slog.String("error", err.Error()))
		}
	}

	// Finished execution
}
