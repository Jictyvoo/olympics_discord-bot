package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type StageRepo struct{ *repoSQLite }

func NewStageRepo(base *repoSQLite) StageRepo { return StageRepo{base} }

func (r StageRepo) UpsertStage(s eventcore.Stage) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertStage(qctx, dbgen.UpsertStageParams{
		ID:          s.ID.Bytes(),
		ProviderID:  s.Ext.Provider,
		ExternalKey: s.Ext.Key,
		Name:        s.Name,
		Ord:         int64(s.Ord),
		SeasonID:    s.SeasonID.Bytes(),
	})
}
