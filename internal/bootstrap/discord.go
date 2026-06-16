package bootstrap

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/usecases/discordsync"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/subscriptions"
	"github.com/jictyvoo/olhojogo/internal/infra/discordfacade"
)

// SetupDiscord builds the discordgo-backed Client from an open session, installs
// the slash-command router, and binds the client into DI for the notifier and
// discordsync consumers. It must run after session.Open(); route-install
// failures are returned (fatal to the caller).
func SetupDiscord(inj remy.Injector, session *discordgo.Session, guildID string) error {
	client := discordfacade.New(session)

	svc, err := remy.Get[subscriptions.Service](inj)
	if err != nil {
		return err
	}

	router := discordfacade.NewRouter(
		"olhojogo", "Manage result notification subscriptions", slog.Default(),
	).
		Add(discordfacade.AddCmd(svc)).
		Add(discordfacade.RemoveCmd(svc)).
		Add(discordfacade.ListCmd(svc))

	if regErr := client.InstallRouter(router, guildID); regErr != nil {
		return regErr
	}

	// Bind the single client instance into the graph for its consumers.
	remy.RegisterInstance[notifier.Dispatcher](inj, client)
	remy.RegisterInstance[notifier.ChannelEnsurer](inj, client)
	remy.RegisterInstance[discordsync.ScheduledEventFacade](inj, client)
	remy.RegisterInstance[notifier.MentionResolver](inj, svc)
	return nil
}
