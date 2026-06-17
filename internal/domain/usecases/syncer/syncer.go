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

// Syncer orchestrates one sync cycle. The ctx is injected at resolution time for
// the provider's outbound HTTP calls. Each persisted fixture is emitted on the
// shared events subject so observers (notifier, discordsync) can react.
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
