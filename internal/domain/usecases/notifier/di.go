package notifier

import (
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
)

// Register wires the Notifier. The factory runs after the Discord session is up,
// so it resolves the channel name to an ID (creating it if missing) first.
func Register(inj remy.Injector, channelName, guildID string, window time.Duration) {
	remy.RegisterFactory(inj, func(ret remy.DependencyRetriever) (*Notifier, error) {
		return New(
			remy.MustGet[FixtureReader](ret),
			remy.MustGet[NotificationRepo](ret),
			remy.MustGet[Dispatcher](ret),
			remy.MustGet[FixtureContextReader](ret),
			remy.MustGet[CompetitorReader](ret),
			render.Olympics{},
			remy.MustGet[MentionResolver](ret),
			ensureChannel(ret, guildID, channelName),
			guildID,
			window,
		), nil
	})
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
