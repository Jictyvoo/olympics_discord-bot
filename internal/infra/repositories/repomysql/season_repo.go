package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type SeasonRepo struct{ *repoMySQL }

func NewSeasonRepo(base *repoMySQL) SeasonRepo { return SeasonRepo{base} }

func (r SeasonRepo) UpsertSeason(s eventcore.Season) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertSeason(qctx, dbgen.UpsertSeasonParams{
		ID:            s.ID.Bytes(),
		ProviderID:    s.Ext.Provider,
		ExternalKey:   s.Ext.Key,
		Name:          s.Name,
		StartsOn:      s.StartsOn.UTC(),
		EndsOn:        s.EndsOn.UTC(),
		CompetitionID: s.CompetitionID.Bytes(),
	})
}
