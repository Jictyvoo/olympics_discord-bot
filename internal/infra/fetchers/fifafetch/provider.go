package fifafetch

import (
	"context"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

type Provider struct{}

func New() Provider { return Provider{} }

func (p Provider) Code() eventcore.ProviderID { return eventcore.ProviderFIFA }
func (p Provider) DisplayName() string        { return "FIFA" }

func (p Provider) SyncFixturesByDate(_ context.Context, _ time.Time) (eventcore.SyncDelta, error) {
	return eventcore.SyncDelta{}, eventcore.ErrNotImplemented
}

func (p Provider) SyncFixtureResults(
	_ context.Context,
	_ eventcore.Fixture,
) (eventcore.SyncDelta, error) {
	return eventcore.SyncDelta{}, eventcore.ErrNotImplemented
}
