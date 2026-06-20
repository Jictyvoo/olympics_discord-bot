package notifier

import "github.com/jictyvoo/olhojogo/internal/domain/eventcore"

// channelRouter resolves the destination channel for a fixture's provider,
// falling back to the default when a provider has no dedicated channel.
type channelRouter struct {
	byProvider map[eventcore.ProviderID]string
	fallback   string
}

func (r channelRouter) channelFor(provider eventcore.ProviderID) string {
	if id := r.byProvider[provider]; id != "" {
		return id
	}
	return r.fallback
}
