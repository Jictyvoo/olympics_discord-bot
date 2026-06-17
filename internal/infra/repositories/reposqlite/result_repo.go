package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type ResultRepo struct{ *repoSQLite }

func NewResultRepo(base *repoSQLite) ResultRepo { return ResultRepo{base} }

func (r ResultRepo) UpsertResult(res eventcore.Result) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	var pos any
	if res.Position != nil {
		pos = int64(*res.Position)
	}
	return r.Queries().UpsertResult(qctx, dbgen.UpsertResultParams{
		FixtureID:     res.FixtureID.Bytes(),
		ParticipantID: res.ParticipantID.Bytes(),
		Position:      pos,
		Score:         mapper.OptString(res.Score),
		RawMark:       mapper.OptString(res.RawMark),
		Outcome:       mapper.OptString(string(res.Outcome)),
	})
}
