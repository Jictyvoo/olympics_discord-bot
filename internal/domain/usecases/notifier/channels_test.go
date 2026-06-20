package notifier

import (
	"testing"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const routerFallback = "default"

func TestChannelRouter_ChannelFor(t *testing.T) {
	router := channelRouter{
		byProvider: map[eventcore.ProviderID]string{eventcore.ProviderVNL: "volei"},
		fallback:   routerFallback,
	}
	testCases := []struct {
		name     string
		provider eventcore.ProviderID
		want     string
	}{
		{"mapped provider uses its channel", eventcore.ProviderVNL, "volei"},
		{"unmapped provider falls back", eventcore.ProviderFIFA, routerFallback},
		{"empty provider falls back", "", routerFallback},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if got := router.channelFor(tCase.provider); got != tCase.want {
				t.Fatalf("channelFor(%q) = %q, want %q", tCase.provider, got, tCase.want)
			}
		})
	}
}

func TestChannelRouter_NilMapFallsBack(t *testing.T) {
	router := channelRouter{fallback: routerFallback}
	if got := router.channelFor(eventcore.ProviderVNL); got != routerFallback {
		t.Fatalf("channelFor on nil map = %q, want default fallback", got)
	}
}
