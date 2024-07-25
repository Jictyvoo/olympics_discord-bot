package reposqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/usecases"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/repositories/reposqlite/internal/dbgen"
)

type RepoSQLite struct {
	conn    *sql.DB
	queries *dbgen.Queries
}

func NewRepoSQLite(db *sql.DB) usecases.AccessDatabaseRepository {
	return RepoSQLite{conn: db, queries: dbgen.New(db)}
}

func (r RepoSQLite) Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
