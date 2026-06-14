package repomysql

import (
	"context"
	"database/sql"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/syncer"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/repobase"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

const hoursPerDay = 24

type repoMySQL = repobase.Base[*dbgen.Queries]

func newRepo(ctx context.Context, db *sql.DB) *repoMySQL {
	return repobase.New[*dbgen.Queries](ctx, db, dbgen.New(db))
}

// Repository is the syncer-facing entry point. SaveCursor / RecordError run on
// the plain db connection through the shared base, never on a finished tx.
type Repository struct {
	*repobase.GenericRepository[*dbgen.Queries, syncer.Tx]
}

func NewRepository(base *repoMySQL) *Repository {
	return &Repository{
		repobase.NewGenericRepository[*dbgen.Queries, syncer.Tx](base, newRepo, newTxAdapter),
	}
}

func (r *Repository) SaveCursor(provider eventcore.ProviderID, scope, cursor string) error {
	return (SyncStateRepo{r.Base}).SaveCursor(provider, scope, cursor)
}

func (r *Repository) RecordError(provider eventcore.ProviderID, scope, errMsg string) error {
	return (SyncStateRepo{r.Base}).RecordError(provider, scope, errMsg)
}

//nolint:ireturn // factory returns the consumer Tx port by design
func newTxAdapter(base *repoMySQL, finish repobase.OnFinishFunc) syncer.Tx {
	return &txAdapter{
		finish:          finish,
		CompetitionRepo: CompetitionRepo{base},
		SeasonRepo:      SeasonRepo{base},
		StageRepo:       StageRepo{base},
		GroupRepo:       GroupRepo{base},
		VenueRepo:       VenueRepo{base},
		ParticipantRepo: ParticipantRepo{base},
		FixtureRepo:     FixtureRepo{base},
		ResultRepo:      ResultRepo{base},
		StandingRepo:    StandingRepo{base},
	}
}

// txAdapter promotes the embedded repos' Upsert* methods; all share one
// tx-bound base, and Commit/Rollback delegate to the finisher.
type txAdapter struct {
	CompetitionRepo
	SeasonRepo
	StageRepo
	GroupRepo
	VenueRepo
	ParticipantRepo
	FixtureRepo
	ResultRepo
	StandingRepo

	finish repobase.OnFinishFunc
}

func (a *txAdapter) Commit() error   { return a.finish(true) }
func (a *txAdapter) Rollback() error { return a.finish(false) }
