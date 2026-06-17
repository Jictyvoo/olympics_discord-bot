package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type ResultRepo struct{ *repoMySQL }

func NewResultRepo(base *repoMySQL) ResultRepo { return ResultRepo{base} }

func (r ResultRepo) UpsertResult(res eventcore.Result) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertResult(qctx, dbgen.UpsertResultParams{
		FixtureID:     res.FixtureID.Bytes(),
		ParticipantID: res.ParticipantID.Bytes(),
		Position:      mapper.NSIntFromPtr(res.Position),
		Score:         mapper.NSStr(res.Score),
		RawMark:       mapper.NSStr(res.RawMark),
		Outcome:       mapper.NSStr(string(res.Outcome)),
	})
}
