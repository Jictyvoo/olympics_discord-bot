package provider

import (
	"context"
	"iter"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

type Strategy interface {
	Code() eventcore.ProviderID
	DisplayName() string
	SyncFixturesByDate(ctx context.Context, day time.Time) (eventcore.SyncDelta, error)
	SyncFixtureResults(ctx context.Context, f eventcore.Fixture) (eventcore.SyncDelta, error)
}

type Set interface {
	Get(p eventcore.ProviderID) (Strategy, error)
	Enabled() iter.Seq[Strategy]
}
