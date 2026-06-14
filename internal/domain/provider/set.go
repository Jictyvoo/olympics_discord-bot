package provider

import (
	"fmt"
	"iter"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

type providerSet struct {
	providers map[eventcore.ProviderID]Strategy
	enabled   []eventcore.ProviderID
}

//nolint:ireturn // factory returning consumer interface by design
func NewSet(all []Strategy, enabledCodes []eventcore.ProviderID) Set {
	m := make(map[eventcore.ProviderID]Strategy, len(all))
	for _, p := range all {
		m[p.Code()] = p
	}
	return &providerSet{providers: m, enabled: enabledCodes}
}

//nolint:ireturn // factory returning consumer interface by design
func (s *providerSet) Get(p eventcore.ProviderID) (Strategy, error) {
	strategy, ok := s.providers[p]
	if !ok {
		return nil, fmt.Errorf("provider: %q not registered", p)
	}
	return strategy, nil
}

func (s *providerSet) Enabled() iter.Seq[Strategy] {
	return func(yield func(Strategy) bool) {
		for _, code := range s.enabled {
			strategy, ok := s.providers[code]
			if !ok {
				continue
			}
			if !yield(strategy) {
				return
			}
		}
	}
}
