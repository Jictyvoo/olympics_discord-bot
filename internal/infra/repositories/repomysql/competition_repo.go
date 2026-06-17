package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type CompetitionRepo struct{ *repoMySQL }

func NewCompetitionRepo(base *repoMySQL) CompetitionRepo { return CompetitionRepo{base} }

// GetFixtureContext resolves the competition, stage and (optional) group that
// locate a fixture in a single query.
func (r CompetitionRepo) GetFixtureContext(
	fixtureID eventcore.CanonicalID,
) (eventcore.FixtureContext, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetFixtureContext(qctx, fixtureID.Bytes())
	if err != nil {
		return eventcore.FixtureContext{}, err
	}
	return eventcore.FixtureContext{
		Competition: eventcore.Competition{
			Code:       row.CompetitionCode.String,
			Name:       row.CompetitionName,
			Discipline: row.Discipline.String,
		},
		StageName: row.StageName,
		StageOrd:  int(row.StageOrd),
		GroupName: row.GroupName.String,
	}, nil
}

func (r CompetitionRepo) UpsertCompetition(c eventcore.Competition) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertCompetition(qctx, dbgen.UpsertCompetitionParams{
		ID:          c.ID.Bytes(),
		ProviderID:  c.Ext.Provider,
		ExternalKey: c.Ext.Key,
		Code:        mapper.NSStr(c.Code),
		Name:        c.Name,
		Discipline:  mapper.NSStr(c.Discipline),
	})
}
