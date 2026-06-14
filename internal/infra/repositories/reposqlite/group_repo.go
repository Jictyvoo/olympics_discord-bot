package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type GroupRepo struct{ *repoSQLite }

func NewGroupRepo(base *repoSQLite) GroupRepo { return GroupRepo{base} }

func (r GroupRepo) UpsertGroup(g eventcore.Group) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertGroup(qctx, dbgen.UpsertGroupParams{
		ID:          g.ID.Bytes(),
		ProviderID:  g.Ext.Provider,
		ExternalKey: g.Ext.Key,
		Name:        g.Name,
		StageID:     g.StageID.Bytes(),
	})
}
