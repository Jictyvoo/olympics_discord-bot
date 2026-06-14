package reposqlite

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/internal/mapper"
	"github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite/dbgen"
)

type DiscordEventRepo struct{ *repoSQLite }

func NewDiscordEventRepo(base *repoSQLite) DiscordEventRepo { return DiscordEventRepo{base} }

func (r DiscordEventRepo) UpsertDiscordEvent(
	de eventcore.DiscordEvent,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpsertDiscordEvent(qctx, dbgen.UpsertDiscordEventParams{
		FixtureID:      de.FixtureID.Bytes(),
		GuildID:        de.GuildID,
		DiscordEventID: de.DiscordEventID,
		Status:         string(de.Status),
		LastChecksum:   mapper.OptString(de.LastChecksum),
	})
}

func (r DiscordEventRepo) GetDiscordEventByFixture(
	fixtureID eventcore.CanonicalID,
	guildID string,
) (eventcore.DiscordEvent, error) {
	qctx, cancel := r.Ctx()
	defer cancel()
	row, err := r.Queries().GetDiscordEventByFixture(qctx, dbgen.GetDiscordEventByFixtureParams{
		FixtureID: fixtureID.Bytes(),
		GuildID:   guildID,
	})
	if err != nil {
		return eventcore.DiscordEvent{}, err
	}
	return rowToDiscordEvent(row), nil
}

func (r DiscordEventRepo) UpdateDiscordEventStatus(
	fixtureID eventcore.CanonicalID,
	guildID string,
	status eventcore.DiscordEventStatus,
) error {
	qctx, cancel := r.Ctx()
	defer cancel()
	return r.Queries().UpdateDiscordEventStatus(qctx, dbgen.UpdateDiscordEventStatusParams{
		Status:    string(status),
		FixtureID: fixtureID.Bytes(),
		GuildID:   guildID,
	})
}

func rowToDiscordEvent(row dbgen.DiscordEvent) eventcore.DiscordEvent {
	return eventcore.DiscordEvent{
		FixtureID:      mapper.IDFromBytes(row.FixtureID),
		GuildID:        row.GuildID,
		DiscordEventID: row.DiscordEventID,
		Status:         eventcore.DiscordEventStatus(row.Status),
		LastChecksum:   mapper.NullStr(row.LastChecksum),
		UpdatedAt:      row.UpdatedAt.UTC(),
	}
}
