package repositories

import (
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/repolympicfetch"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite"

	"github.com/wrapped-owls/goremy-di/remy"
	_ "modernc.org/sqlite"
)

func RegisterRepositories(inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[usecases.AccessDatabaseRepository], reposqlite.NewRepoSQLite,
	)
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[usecases.OlympicsFetcher], repolympicfetch.NewOlympicsFetcher,
	)
}
