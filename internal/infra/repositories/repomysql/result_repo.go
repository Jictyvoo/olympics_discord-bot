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

func (r ResultRepo) ListResultsByFixture(
	fixtureID eventcore.CanonicalID,
) ([]eventcore.Result, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListResultsByFixture(qctx, fixtureID.Bytes())
	if err != nil {
		return nil, err
	}
	out := make([]eventcore.Result, len(rows))
	for i, row := range rows {
		out[i] = eventcore.Result{
			FixtureID:     mapper.IDFromBytes(row.FixtureID),
			ParticipantID: mapper.IDFromBytes(row.ParticipantID),
			Position:      mapper.IntFromNull64(row.Position),
			Score:         row.Score.String,
			RawMark:       row.RawMark.String,
			Outcome:       eventcore.Outcome(row.Outcome.String),
		}
	}
	return out, nil
}
