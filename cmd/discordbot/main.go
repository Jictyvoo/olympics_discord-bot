package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olympics_data_fetcher/internal/bootstrap"
	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/services"
	"github.com/jictyvoo/olympics_data_fetcher/pkg/config"
)

func main() {
	conf, confErr := config.LoadTOML(config.DefaultFileName)
	if confErr != nil {
		slog.Error(
			"Error loading config file",
			slog.String("file", config.DefaultFileName),
			slog.String("error", confErr.Error()),
		)
		os.Exit(1)
	}

	config.LoadConfigFromLoader(&conf, config.EnvLoader{})
	db := bootstrap.OpenDatabase()
	defer db.Close()

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance(inj, db)
	bootstrap.DoInjections(inj)

	fmt.Println(generateInviteLink(conf.Discord.ClientID))
	discClient, err := discordgo.New("Bot " + conf.Discord.Token)
	if err != nil {
		slog.Error(
			"Error creating Discord session",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	cancelChan := make(services.CancelChannel, 1)
	notifierServ, err := remy.DoGetGenFunc[*services.EventNotifier](
		inj, func(injector remy.Injector) error {
			remy.RegisterInstance(injector, cancelChan)
			return nil
		},
	)
	if err != nil {
		slog.Error(
			"Error getting notifier service",
			slog.String("error", err.Error()),
		)
		return
	}
	bot := NewOlympicsBot(notifierServ)
	discClient.AddHandler(bot.ReadyHandler)
	discClient.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent,
	)

	gracefullShutdown(cancelChan, notifierServ, discClient)
}
