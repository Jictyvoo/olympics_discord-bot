package repobase

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const defaultQueryTimeoutSeconds = 30

// Queries is the sqlc query set, rebindable to a tx. Each dialect's
// *dbgen.Queries satisfies it with Q set to its own pointer type.
type Queries[Q any] interface {
	WithTx(tx *sql.Tx) Q
}

// OnFinishFunc commits (true) or rolls back (false) a transaction.
type OnFinishFunc func(commit bool) error

// Base is embedded (by pointer) into every concrete repo. BeginTx swaps queries
// to a tx-bound clone, so all repos sharing one Base write through a single
// transaction. The injected ctx drives every per-query timeout through Ctx.
type Base[Q Queries[Q]] struct {
	ctx     context.Context
	conn    *sql.DB
	queries Q
}

func New[Q Queries[Q]](ctx context.Context, conn *sql.DB, queries Q) *Base[Q] {
	return &Base[Q]{ctx: ctx, conn: conn, queries: queries}
}

func (b *Base[Q]) Queries() Q               { return b.queries }
func (b *Base[Q]) Connection() *sql.DB      { return b.conn }
func (b *Base[Q]) Context() context.Context { return b.ctx }

func (b *Base[Q]) BeginTx(ctx context.Context, txOpts *sql.TxOptions) (OnFinishFunc, error) {
	tx, err := b.conn.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}
	b.queries = b.queries.WithTx(tx)
	return func(commit bool) error {
		if commit {
			return tx.Commit()
		}
		// Rollback after Commit is a safe no-op: database/sql returns sql.ErrTxDone.
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			return rbErr
		}
		return nil
	}, nil
}

// Ctx derives a per-query context with the default timeout from the injected
// base context, falling back to Background when none was injected.
func (b *Base[Q]) Ctx() (context.Context, context.CancelFunc) {
	parent := b.ctx
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, defaultQueryTimeoutSeconds*time.Second)
}

// GenericRepository implements the syncer-facing Begin once for all dialects;
// each dialect supplies the newBase / newTx constructors.
type GenericRepository[Q Queries[Q], T any] struct {
	*Base[Q]

	newBase func(context.Context, *sql.DB) *Base[Q]
	newTx   func(*Base[Q], OnFinishFunc) T
}

func NewGenericRepository[Q Queries[Q], T any](
	base *Base[Q],
	newBase func(context.Context, *sql.DB) *Base[Q],
	newTx func(*Base[Q], OnFinishFunc) T,
) *GenericRepository[Q, T] {
	return &GenericRepository[Q, T]{Base: base, newBase: newBase, newTx: newTx}
}

func (r *GenericRepository[Q, T]) Begin() (T, error) {
	base := r.newBase(r.Context(), r.Connection())
	bctx, cancel := base.Ctx()
	finish, err := base.BeginTx(bctx, nil)
	if err != nil {
		cancel()
		var zero T
		return zero, err
	}
	// Release the tx context only once the transaction finishes; cancelling at
	// Begin return would make database/sql roll the tx back out from under the
	// caller before it can Commit.
	return r.newTx(base, func(commit bool) error {
		defer cancel()
		return finish(commit)
	}), nil
}
