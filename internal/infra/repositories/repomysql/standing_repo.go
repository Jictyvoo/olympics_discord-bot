package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type StandingRepo struct{ *repoMySQL }

func NewStandingRepo(base *repoMySQL) StandingRepo { return StandingRepo{base} }

// UpsertStanding persists a stage standing. The Stats map has no dedicated
// column in the generic schema and is intentionally not persisted here.
func (r StandingRepo) UpsertStanding(s eventcore.Standing) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertStanding(qctx, dbgen.UpsertStandingParams{
		StageID:       s.StageID.Bytes(),
		ParticipantID: s.ParticipantID.Bytes(),
		Rank:          int64(s.Rank),
		Points:        int64(s.Points),
	})
}
