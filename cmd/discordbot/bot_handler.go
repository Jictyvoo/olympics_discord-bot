package main

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/services"
	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/facades/discfac"
)

const channelName = "🏅olympics_PARIS-2024"

type OlympicsBot struct {
	evtNotifier    *services.EventNotifier
	watchCountries []string
	guildIDs       map[string]services.NotifierFacade
}

func NewOlympicsBot(evtNotifier *services.EventNotifier, watchCountries []string) *OlympicsBot {
	return &OlympicsBot{
		evtNotifier:    evtNotifier,
		watchCountries: watchCountries,
		guildIDs:       make(map[string]services.NotifierFacade, 11),
	}
}

func (bot *OlympicsBot) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	for _, guild := range r.Guilds {
		if _, ok := bot.guildIDs[guild.ID]; ok {
			continue
		}
		discFacade := discfac.NewDiscordFacadeImpl(guild.ID, s)
		if err := discFacade.InitMessageChannel(channelName); err != nil {
			slog.Error(
				"Error initializing discord channel",
				slog.String("guild", guild.ID),
				slog.String("err", err.Error()),
			)
		}
		bot.guildIDs[guild.ID] = discFacade

		eventManager, err := services.NewOlympicEventManager(bot.watchCountries, discFacade)
		if err != nil {
			slog.Error(
				"Error initializing event manager",
				slog.String("guild", guild.ID),
				slog.String("err", err.Error()),
			)
			return
		}
		// Register event observer
		bot.evtNotifier.RegisterObserver(eventManager)
	}

	// Start scheduler
	bot.evtNotifier.Start()
}

func (bot *OlympicsBot) MessagesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	/*_, err := s.ChannelMessageSend(
		m.ChannelID, "Future feature, isn't implemented yet.",
	)
	if err != nil {
		slog.Error(
			"Error sending message",
			slog.String("err", err.Error()),
		)
	}*/
}
