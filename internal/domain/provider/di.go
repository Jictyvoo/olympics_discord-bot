package provider

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// RegisterSet builds the provider Set. Each fetcher must already be registered
// as a Strategy keyed by its ProviderID; bootstrap orchestrates that order.
func RegisterSet(inj remy.Injector, enabled []eventcore.ProviderID) {
	all := make([]Strategy, 0, len(enabled))
	for _, code := range enabled {
		s, err := remy.Get[Strategy](inj, code)
		if err != nil {
			continue
		}
		all = append(all, s)
	}
	set := NewSet(all, enabled)
	remy.RegisterInstance(inj, set)
}
