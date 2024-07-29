package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/jictyvoo/olympics_data_fetcher/internal/domain/services"
)

func gracefullShutdown(
	cancelChan services.CancelChannel, notifierServ *services.EventNotifier,
	discClient *discordgo.Session,
) {
	var wg sync.WaitGroup

	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := notifierServ.MainLoop(); err != nil {
				slog.Error(
					"Error starting notifier service",
					slog.String("err", err.Error()),
				)
				return
			}
		}()
	}

	// Open a websocket connection to Discord and begin listening.
	{
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := discClient.Open(); err != nil {
				slog.Error(
					"Error opening connection",
					slog.String("err", err.Error()),
				)
				return
			}
		}()
	}

	// Cleanly close down the Discord session.
	defer discClient.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	cancelChan <- struct{}{}
	wg.Wait()
}
