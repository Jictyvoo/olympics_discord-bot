package syncer

import (
	"context"
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"
)

const defaultSyncIntervalMinutes = 4

// syncLookbackDays is how many days before today each tick re-syncs, so a
// fixture whose feed finalizes after its day rolls over still gets picked up.
const syncLookbackDays = 1

// Runner drives the sync loop: ticks every interval and resolves a Syncer bound
// to the tick's context (remy.GetWithContext), then calls SyncDay for today.
// Resolving per tick is what lets each cycle's repository + provider calls honour
// the tick context (including cancellation on shutdown).
type Runner struct {
	inj      remy.DependencyRetriever
	interval time.Duration
}

func NewRunner(inj remy.DependencyRetriever, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = defaultSyncIntervalMinutes * time.Minute
	}
	return &Runner{inj: inj, interval: interval}
}

// Run blocks until ctx is cancelled, calling SyncDay on each tick.
func (r *Runner) Run(ctx context.Context) error {
	// Sync once immediately before first tick.
	r.tick(ctx)

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			r.tick(ctx)
		}
	}
}

func (r *Runner) tick(ctx context.Context) {
	now := time.Now().UTC()
	dailySyncer, err := remy.GetWithContext[*Syncer](r.inj, ctx)
	if err != nil {
		slog.Error("syncer: resolve syncer", slog.String("err", err.Error()))
		return
	}
	from := now.AddDate(0, 0, -syncLookbackDays)
	if err = dailySyncer.SyncRange(from, now); err != nil {
		slog.Error(
			"syncer: SyncRange failed",
			slog.String("err", err.Error()),
			slog.Time("from", from),
			slog.Time("to", now),
		)
	}
}
