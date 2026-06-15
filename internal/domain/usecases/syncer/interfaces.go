package syncer

//go:generate go tool -modfile=../../../../tools/go.mod mockgen -source=interfaces.go -destination=interfaces_mock_test.go -package=syncer

import (
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// Repository opens persistence transactions and records per-provider sync state.
// The context is bound at resolution time (remy.GetWithContext), so these
// methods no longer take a context.Context.
type Repository interface {
	Begin() (Tx, error)
	SaveCursor(provider eventcore.ProviderID, scope, cursor string) error
	RecordError(provider eventcore.ProviderID, scope, errMsg string) error
}

// Tx upserts every record of a sync delta within a single transaction.
//
//nolint:interfacebloat // one upsert per persisted entity type, by design
type Tx interface {
	UpsertCompetition(c eventcore.Competition) error
	UpsertSeason(s eventcore.Season) error
	UpsertStage(s eventcore.Stage) error
	UpsertGroup(g eventcore.Group) error
	UpsertVenue(v eventcore.Venue) error
	UpsertParticipant(p eventcore.Participant) error
	UpsertFixture(f eventcore.Fixture) error
	UpsertFixtureParticipants(
		fixtureID eventcore.CanonicalID,
		parts []eventcore.FixtureParticipant,
	) error
	UpsertResult(res eventcore.Result) error
	UpsertStanding(s eventcore.Standing) error
	Commit() error
	Rollback() error
}
