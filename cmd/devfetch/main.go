package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/bootstrap"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	_ "github.com/jictyvoo/olhojogo/internal/infra/repositories/dbdrivers"
)

const exitUsage = 2

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if len(os.Args) < exitUsage {
		slog.Error("usage: devfetch <sync|version>")
		os.Exit(exitUsage)
	}

	switch os.Args[1] {
	case "sync":
		opts, err := parseSyncArgs(os.Args[2:])
		if err != nil {
			slog.Error("sync args", slog.String("err", err.Error()))
			os.Exit(exitUsage)
		}
		if err := cmdSync(opts); err != nil {
			slog.Error("sync failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	case "version":
		if _, err := os.Stdout.WriteString("devfetch dev\n"); err != nil {
			slog.Error("write version", slog.String("err", err.Error()))
		}
	default:
		slog.Error("unknown subcommand", slog.String("cmd", os.Args[1]))
		os.Exit(exitUsage)
	}
}

func cmdSync(opts syncOptions) error {
	conf, err := appconfig.Load("")
	if err != nil {
		return err
	}

	db, err := sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	bootstrap.DoInjections(inj, conf, db)

	providers, err := remy.Get[provider.Set](inj)
	if err != nil {
		return err
	}
	strategy, err := providers.Get(opts.Provider)
	if err != nil {
		return err
	}

	delta, err := strategy.SyncFixturesByDate(context.Background(), opts.Date)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(delta)
}
