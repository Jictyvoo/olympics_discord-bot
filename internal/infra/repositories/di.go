package repositories

import (
	"database/sql"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/services"
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/repolympicfetch"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite"

	"github.com/wrapped-owls/goremy-di/remy"
	_ "modernc.org/sqlite"
)

func constructRepoSQLITE[T any](db *sql.DB) (value T) {
	var result any
	result = reposqlite.NewRepoSQLite(db)

	if asT, ok := result.(T); ok {
		return asT
	}

	return
}

func RegisterRepositories(inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj,
		remy.Factory[usecases.AccessDatabaseRepository],
		constructRepoSQLITE[usecases.AccessDatabaseRepository],
	)
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[services.EventLoader], constructRepoSQLITE[services.EventLoader],
	)
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[usecases.OlympicsFetcher], repolympicfetch.NewOlympicsFetcher,
	)
}
