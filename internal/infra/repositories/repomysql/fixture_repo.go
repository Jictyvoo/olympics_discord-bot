package repomysql

import (
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/repomysql/dbgen"
)

type FixtureRepo struct{ *repoMySQL }

func NewFixtureRepo(base *repoMySQL) FixtureRepo { return FixtureRepo{base} }

func (r FixtureRepo) UpsertFixture(f eventcore.Fixture) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertFixture(qctx, dbgen.UpsertFixtureParams{
		ID:          f.ID.Bytes(),
		ProviderID:  f.Ext.Provider,
		ExternalKey: f.Ext.Key,
		StageID:     f.StageID.Bytes(),
		GroupID:     mapper.NSStrFromID(f.GroupID),
		VenueID:     mapper.NSStrFromID(f.VenueID),
		Name:        f.Name,
		StartsAt:    f.StartsAt.UTC(),
		EndsAt:      f.EndsAt.UTC(),
		Status:      string(f.Status),
		Checksum:    mapper.NSStr(f.Checksum),
	})
}

func (r FixtureRepo) GetFixture(
	id eventcore.CanonicalID,
) (eventcore.Fixture, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetFixture(qctx, id.Bytes())
	if err != nil {
		return eventcore.Fixture{}, err
	}
	return rowToFixture(row), nil
}

func (r FixtureRepo) ListFixturesByDay(
	provider eventcore.ProviderID,
	day time.Time,
) ([]eventcore.Fixture, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	start := day.UTC().Truncate(hoursPerDay * time.Hour)
	rows, err := r.Queries().ListFixturesByDay(qctx, dbgen.ListFixturesByDayParams{
		ProviderID: provider,
		StartsAt:   start,
		StartsAt_2: start.Add(hoursPerDay * time.Hour),
	})
	if err != nil {
		return nil, err
	}
	return rowsToFixtures(rows), nil
}

func (r FixtureRepo) ListFixturesStartingBefore(
	before time.Time,
) ([]eventcore.Fixture, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	rows, err := r.Queries().ListFixturesStartingBefore(qctx, before.UTC())
	if err != nil {
		return nil, err
	}
	return rowsToFixtures(rows), nil
}

func (r FixtureRepo) UpdateFixtureStatus(
	id eventcore.CanonicalID,
	status eventcore.FixtureStatus,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpdateFixtureStatus(qctx, dbgen.UpdateFixtureStatusParams{
		Status: string(status),
		ID:     id.Bytes(),
	})
}

func rowToFixture(row dbgen.Fixture) eventcore.Fixture {
	return eventcore.Fixture{
		ID:       mapper.IDFromBytes(row.ID),
		Ext:      eventcore.ExternalID{Provider: row.ProviderID, Key: row.ExternalKey},
		StageID:  mapper.IDFromBytes(row.StageID),
		GroupID:  mapper.IDFromNullStr(row.GroupID),
		VenueID:  mapper.IDFromNullStr(row.VenueID),
		Name:     row.Name,
		StartsAt: row.StartsAt.UTC(),
		EndsAt:   row.EndsAt.UTC(),
		Status:   eventcore.FixtureStatus(row.Status),
		Checksum: row.Checksum.String,
	}
}

func rowsToFixtures(rows []dbgen.Fixture) []eventcore.Fixture {
	out := make([]eventcore.Fixture, len(rows))
	for i, row := range rows {
		out[i] = rowToFixture(row)
	}
	return out
}
