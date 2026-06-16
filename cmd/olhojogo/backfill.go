package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/bootstrap"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/syncer"
	_ "github.com/jictyvoo/olhojogo/internal/infra/repositories/dbdrivers"
)

// backfill syncs and persists every day in the inclusive [from, to] range.
func backfill(opts cliOptions) error {
	from, to, err := backfillRange(opts.From, opts.To)
	if err != nil {
		return err
	}

	conf, err := appconfig.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if conf.Database.RunMigrations {
		if migrErr := runMigrations(conf); migrErr != nil {
			return migrErr
		}
	}

	db, err := sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := signalContext()
	defer cancel()

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	bootstrap.DoInjections(inj, conf, db)

	dailySyncer, err := remy.GetWithContext[*syncer.Syncer](inj, ctx)
	if err != nil {
		return err
	}

	slog.Info(
		"backfill started",
		slog.String("from", from.Format(time.DateOnly)),
		slog.String("to", to.Format(time.DateOnly)),
	)
	if err = dailySyncer.SyncRange(from, to); err != nil {
		// Per-day errors are joined; each day still persisted what it could.
		slog.Warn("backfill completed with errors", slog.String("err", err.Error()))
	} else {
		slog.Info("backfill completed")
	}
	return nil
}

// backfillRange requires -from; -to defaults to today (UTC) when empty.
func backfillRange(fromStr, toStr string) (from, to time.Time, err error) {
	if fromStr == "" {
		return from, to, fmt.Errorf("backfill: -from is required (YYYY-MM-DD)")
	}
	from, err = time.Parse(time.DateOnly, fromStr)
	if err != nil {
		return from, to, fmt.Errorf("backfill: parse -from: %w", err)
	}
	if toStr == "" {
		to = time.Now().UTC()
	} else if to, err = time.Parse(time.DateOnly, toStr); err != nil {
		return from, to, fmt.Errorf("backfill: parse -to: %w", err)
	}
	if to.Before(from) {
		return from, to, fmt.Errorf("backfill: -to %q is before -from %q", toStr, fromStr)
	}
	return from, to, nil
}
