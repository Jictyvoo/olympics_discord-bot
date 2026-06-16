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

// SetupDiscord installs the slash-command router and binds the client into DI
// for its consumers. It must run after session.Open().
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
		Add(discordfacade.ListCmd(svc)).
		Add(discordfacade.CountriesCmd(svc))

	if regErr := client.InstallRouter(router, guildID); regErr != nil {
		return regErr
	}

	remy.RegisterInstance[notifier.Dispatcher](inj, client)
	remy.RegisterInstance[notifier.ChannelEnsurer](inj, client)
	remy.RegisterInstance[discordsync.ScheduledEventFacade](inj, client)
	remy.RegisterInstance[notifier.MentionResolver](inj, svc)
	return nil
}
