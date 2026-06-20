package notifier

import (
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
)

// Register wires the Notifier. The factory runs after the Discord session is up,
// so it resolves channel names to IDs (creating them if missing) first.
// perProvider names override the default for fixtures of that provider.
func Register(
	inj remy.Injector,
	defaultChannel string,
	perProvider map[eventcore.ProviderID]string,
	guildID string,
	window time.Duration,
) {
	remy.RegisterFactory(inj, func(ret remy.DependencyRetriever) (*Notifier, error) {
		return New(
			remy.MustGet[FixtureReader](ret),
			remy.MustGet[NotificationRepo](ret),
			remy.MustGet[Dispatcher](ret),
			remy.MustGet[FixtureContextReader](ret),
			remy.MustGet[CompetitorReader](ret),
			render.Olympics{},
			remy.MustGet[MentionResolver](ret),
			resolveChannels(ret, guildID, defaultChannel, perProvider),
			guildID,
			window,
		), nil
	})
}

func resolveChannels(
	ret remy.DependencyRetriever,
	guildID, defaultChannel string,
	perProvider map[eventcore.ProviderID]string,
) channelRouter {
	router := channelRouter{fallback: ensureChannel(ret, guildID, defaultChannel)}
	if len(perProvider) == 0 {
		return router
	}
	router.byProvider = make(map[eventcore.ProviderID]string, len(perProvider))
	for provider, name := range perProvider {
		router.byProvider[provider] = ensureChannel(ret, guildID, name)
	}
	return router
}

// ensureChannel resolves channelName to a channel ID, falling back to the name.
func ensureChannel(ret remy.DependencyRetriever, guildID, channelName string) string {
	if channelName == "" || guildID == "" {
		return channelName
	}
	ensurer, err := remy.Get[ChannelEnsurer](ret)
	if err != nil {
		return channelName
	}
	id, err := ensurer.ResolveChannel(guildID, channelName)
	if err != nil {
		slog.Warn(
			"notifier: resolve channel",
			slog.String("channel", channelName),
			slog.String("err", err.Error()),
		)
		return channelName
	}
	return id
}
