package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type ParticipantRepo struct{ *repoSQLite }

func NewParticipantRepo(base *repoSQLite) ParticipantRepo { return ParticipantRepo{base} }

func (r ParticipantRepo) UpsertParticipant(p eventcore.Participant) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertParticipant(qctx, dbgen.UpsertParticipantParams{
		ID:          p.ID.Bytes(),
		ProviderID:  p.Ext.Provider,
		ExternalKey: p.Ext.Key,
		Kind:        string(p.Kind),
		Name:        p.Name,
		Code:        mapper.OptString(p.Code),
		CountryIso:  mapper.OptString(p.CountryISO),
		Gender:      mapper.OptString(p.Gender),
	})
}

func (r ParticipantRepo) UpsertFixtureParticipants(
	fixtureID eventcore.CanonicalID,
	parts []eventcore.FixtureParticipant,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	fxBytes := fixtureID.Bytes()
	for _, fp := range parts {
		if err := r.Queries().UpsertFixtureParticipant(qctx, dbgen.UpsertFixtureParticipantParams{
			FixtureID:     fxBytes,
			ParticipantID: fp.ParticipantID.Bytes(),
			Role:          mapper.OptString(fp.Role),
		}); err != nil {
			return err
		}
	}
	return nil
}

// ListFixtureCompetitors returns each participant in a fixture with its role,
// the ISO2 code resolved from its country, and its result (if any) in one query.
func (r ParticipantRepo) ListFixtureCompetitors(
	fixtureID eventcore.CanonicalID,
) ([]eventcore.FixtureCompetitor, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListFixtureCompetitors(qctx, fixtureID.Bytes())
	if err != nil {
		return nil, err
	}
	out := make([]eventcore.FixtureCompetitor, len(rows))
	for i, row := range rows {
		pid := mapper.IDFromBytes(row.ID)
		out[i] = eventcore.FixtureCompetitor{
			Participant: eventcore.Participant{
				ID:         pid,
				Ext:        eventcore.ExternalID{Provider: row.ProviderID, Key: row.ExternalKey},
				Kind:       eventcore.ParticipantKind(row.Kind),
				Name:       row.Name,
				Code:       mapper.NullStr(row.Code),
				CountryISO: mapper.NullStr(row.CountryIso),
				Gender:     mapper.NullStr(row.Gender),
			},
			Role:        mapper.NullStr(row.Role),
			CountryISO2: mapper.NullStr(row.CountryIso2),
			Result: eventcore.Result{
				FixtureID:     fixtureID,
				ParticipantID: pid,
				Position:      mapper.NullInt(row.Position),
				Score:         mapper.NullStr(row.Score),
				Outcome:       eventcore.Outcome(mapper.NullStr(row.Outcome)),
			},
		}
	}
	return out, nil
}
