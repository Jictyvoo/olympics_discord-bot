package reposqlite

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type SyncStateRepo struct{ *repoSQLite }

func NewSyncStateRepo(base *repoSQLite) SyncStateRepo { return SyncStateRepo{base} }

type SyncState struct {
	ProviderID   eventcore.ProviderID
	Scope        string
	Cursor       string
	LastSyncedAt time.Time
	LastError    string
}

func (r SyncStateRepo) SaveCursor(
	provider eventcore.ProviderID,
	scope, cursor string,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertSyncState(qctx, dbgen.UpsertSyncStateParams{
		ProviderID: provider,
		Scope:      scope,
		Cursor:     mapper.OptString(cursor),
	})
}

func (r SyncStateRepo) GetSyncState(
	provider eventcore.ProviderID,
	scope string,
) (SyncState, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetSyncState(qctx, dbgen.GetSyncStateParams{
		ProviderID: provider,
		Scope:      scope,
	})
	if err != nil {
		return SyncState{}, err
	}
	return SyncState{
		ProviderID:   row.ProviderID,
		Scope:        row.Scope,
		Cursor:       mapper.NullStr(row.Cursor),
		LastSyncedAt: mapper.TimeOrZero(row.LastSyncedAt),
		LastError:    mapper.NullStr(row.LastError),
	}, nil
}

func (r SyncStateRepo) RecordError(
	provider eventcore.ProviderID,
	scope, errMsg string,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().SetSyncStateError(qctx, dbgen.SetSyncStateErrorParams{
		LastError:  errMsg,
		ProviderID: provider,
		Scope:      scope,
	})
}
