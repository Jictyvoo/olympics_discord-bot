package repomysql

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type CompetitionRepo struct{ *repoMySQL }

func NewCompetitionRepo(base *repoMySQL) CompetitionRepo { return CompetitionRepo{base} }

// GetCompetitionByFixture resolves the competition that owns a fixture by
// walking fixture -> stage -> season -> competition.
func (r CompetitionRepo) GetCompetitionByFixture(
	fixtureID eventcore.CanonicalID,
) (eventcore.Competition, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetCompetitionByFixture(qctx, fixtureID.Bytes())
	if err != nil {
		return eventcore.Competition{}, err
	}
	return eventcore.Competition{
		ID:         mapper.IDFromBytes(row.ID),
		Ext:        eventcore.ExternalID{Provider: row.ProviderID, Key: row.ExternalKey},
		Code:       row.Code.String,
		Name:       row.Name,
		Discipline: row.Discipline.String,
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
