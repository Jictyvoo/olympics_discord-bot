package syncer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
)

// Syncer orchestrates one sync cycle for a single provider + date. The ctx is
// injected at resolution time and used for the provider's outbound HTTP calls;
// the repository binds its own ctx (remy.GetWithContext) so its methods take none.
// After a fixture is persisted, it is emitted on the shared events subject so
// observers (notifier, discordsync) can react.
type Syncer struct {
	providers provider.Set
	repo      Repository
	events    *services.Subject[eventcore.Fixture]
	ctx       context.Context
}

func New(
	providers provider.Set,
	repo Repository,
	events *services.Subject[eventcore.Fixture],
	ctx context.Context,
) *Syncer {
	return &Syncer{providers: providers, repo: repo, events: events, ctx: ctx}
}

// SyncDay fetches and persists all fixtures for every enabled provider on the given day.
func (s *Syncer) SyncDay(day time.Time) error {
	var errs []error
	for p := range s.providers.Enabled() {
		if err := s.syncProviderDay(p, day); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s *Syncer) syncProviderDay(p provider.Strategy, day time.Time) error {
	scope := day.UTC().Format(time.DateOnly)
	delta, err := p.SyncFixturesByDate(s.ctx, day)
	if err != nil {
		if recErr := s.repo.RecordError(p.Code(), scope, err.Error()); recErr != nil {
			slog.Error("syncer: record error failed", slog.String("err", recErr.Error()))
		}
		return err
	}

	if persistErr := s.persist(delta); persistErr != nil {
		_ = s.repo.RecordError(p.Code(), scope, persistErr.Error())
		return persistErr
	}

	return s.repo.SaveCursor(p.Code(), scope, delta.Cursor)
}

// persist writes every record in the delta within a single transaction in
// FK-safe order so that parents exist before their children:
// Competitions -> Seasons -> Stages -> Groups -> Venues -> Participants ->
// Fixtures (+ FixtureParticipants) -> Results -> Standings.
func (s *Syncer) persist(delta eventcore.SyncDelta) error {
	tx, err := s.repo.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err = persistEach(delta.Competitions, tx.UpsertCompetition); err != nil {
		return err
	}
	if err = persistEach(delta.Seasons, tx.UpsertSeason); err != nil {
		return err
	}
	if err = persistEach(delta.Stages, tx.UpsertStage); err != nil {
		return err
	}
	if err = persistEach(delta.Groups, tx.UpsertGroup); err != nil {
		return err
	}
	if err = persistEach(delta.Venues, tx.UpsertVenue); err != nil {
		return err
	}
	if err = persistEach(delta.Participants, tx.UpsertParticipant); err != nil {
		return err
	}
	if err = persistFixtures(tx, delta.Fixtures); err != nil {
		return err
	}
	if err = persistEach(delta.Results, tx.UpsertResult); err != nil {
		return err
	}
	if err = persistEach(delta.Standings, tx.UpsertStanding); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	// Emit each persisted fixture so observers (notifier, discordsync) react
	// only to durably-committed state.
	if s.events != nil {
		for _, f := range delta.Fixtures {
			s.events.Emit(f)
		}
	}
	return nil
}

// persistFixtures upserts each fixture together with its participant links.
func persistFixtures(tx Tx, fixtures []eventcore.Fixture) error {
	for _, f := range fixtures {
		if err := tx.UpsertFixture(f); err != nil {
			return err
		}
		if err := tx.UpsertFixtureParticipants(f.ID, f.Participants); err != nil {
			return err
		}
	}
	return nil
}

func persistEach[T any](
	items []T,
	upsert func(T) error,
) error {
	for _, item := range items {
		if err := upsert(item); err != nil {
			return err
		}
	}
	return nil
}
