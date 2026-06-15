package syncer

import (
	"context"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
)

func Register(inj remy.Injector, syncInterval time.Duration) {
	// Syncer takes a context.Context arg so remy.GetWithContext flows the
	// per-tick ctx into it (and onward to the provider's outbound HTTP calls).
	// The events subject is a long-lived singleton shared across ticks.
	remy.RegisterFactory(inj, func(ret remy.DependencyRetriever) (*Syncer, error) {
		return New(
			remy.MustGet[provider.Set](ret),
			remy.MustGet[Repository](ret),
			remy.MustGet[*services.Subject[eventcore.Fixture]](ret),
			remy.MustGet[context.Context](ret),
		), nil
	})
	// Runner holds the retriever so it can resolve a ctx-bound Syncer per tick.
	remy.RegisterFactory(inj, func(ret remy.DependencyRetriever) (*Runner, error) {
		return NewRunner(ret, syncInterval), nil
	})
}
