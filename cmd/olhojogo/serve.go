package main

import (
	"database/sql"
	"log/slog"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/bootstrap"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/syncer"
	_ "github.com/jictyvoo/olhojogo/internal/infra/repositories/dbdrivers"
)

func serve(configPath string) error {
	conf, err := appconfig.Load(configPath)
	if err != nil {
		return err
	}

	db, err := sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if conf.Database.RunMigrations {
		if migrErr := runMigrations(conf); migrErr != nil {
			return migrErr
		}
	}

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	bootstrap.DoInjections(inj, conf, db)

	ctx, cancel := signalContext()
	defer cancel()

	var session *discordgo.Session
	if conf.Discord.Token != "" {
		session, err = discordgo.New("Bot " + conf.Discord.Token)
		if err != nil {
			return err
		}
		if openErr := session.Open(); openErr != nil {
			return openErr
		}
		if setupErr := bootstrap.SetupDiscord(inj, session, conf.Discord.GuildID); setupErr != nil {
			_ = session.Close()
			return setupErr
		}
		bootstrap.WireObservers(inj)
	}

	runner, err := remy.Get[*syncer.Runner](inj)
	if err != nil {
		slog.Error("serve: get syncer runner", slog.String("err", err.Error()))
		return err
	}

	slog.Info("olhojogo started", slog.String("project", conf.ProjectName))

	// Close the Discord session once ctx is cancelled, alongside the syncer loop.
	var wg sync.WaitGroup
	var runErr error
	wg.Go(func() { runErr = runner.Run(ctx) })
	if session != nil {
		wg.Go(func() {
			<-ctx.Done()
			_ = session.Close()
		})
	}
	wg.Wait()
	return runErr
}
