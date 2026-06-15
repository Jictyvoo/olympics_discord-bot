package notifier

import (
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
)

// Register wires the Notifier into the DI container. channelID, guildID, and
// window are the runtime parameters; the readers, NotificationRepo, and
// Dispatcher are resolved from the graph. The dependency count exceeds
// RegisterConstructorArgsN, so a plain factory resolves each consumer interface
// from the retriever.
func Register(inj remy.Injector, channelID, guildID string, window time.Duration) {
	remy.RegisterFactory(inj, func(ret remy.DependencyRetriever) (*Notifier, error) {
		return New(
			remy.MustGet[FixtureReader](ret),
			remy.MustGet[NotificationRepo](ret),
			remy.MustGet[Dispatcher](ret),
			remy.MustGet[ResultReader](ret),
			remy.MustGet[CompetitionReader](ret),
			remy.MustGet[ParticipantReader](ret),
			render.Olympics{},
			remy.MustGet[MentionResolver](ret),
			channelID,
			guildID,
			window,
		), nil
	})
}
