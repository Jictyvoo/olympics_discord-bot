package main

import (
	"context"
	"database/sql"
	"log/slog"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/database/migrations/sqlite"
	_ "github.com/jictyvoo/olhojogo/internal/infra/repositories/dbdrivers"
	"github.com/jictyvoo/olhojogo/internal/migrator"
)

func migrate(configPath string) error {
	conf, err := appconfig.Load(configPath)
	if err != nil {
		return err
	}
	return runMigrations(conf)
}

func runMigrations(conf appconfig.Config) error {
	db, err := sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	m := migrator.New(db, conf.Database.Driver, sqlite.FS)
	if migrErr := m.Run(context.Background()); migrErr != nil {
		slog.Error("migration failed", slog.String("err", migrErr.Error()))
		return migrErr
	}
	slog.Info("migrations applied successfully")
	return nil
}
