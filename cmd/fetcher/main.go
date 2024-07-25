package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain"
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra"

	"github.com/wrapped-owls/goremy-di/remy"
)

func main() {
	db, dbErr := sql.Open("sqlite", "olympics-2024_PARIS.db")
	if dbErr != nil {
		slog.Error("failed to open database", slog.String("error", dbErr.Error()))
		os.Exit(1)
	}
	defer db.Close()

	inj := remy.NewInjector(remy.Config{})
	remy.RegisterInstance(inj, ".rest_cache", "cacheDirectory")
	remy.RegisterInstance(inj, db)
	infra.RegisterInfraServices(inj)
	domain.RegisterUCs(inj)

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
