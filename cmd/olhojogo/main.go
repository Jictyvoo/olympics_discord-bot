package main

import (
	"log/slog"
	"os"
)

const exitUsage = 2

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	opts := parseArgs(os.Args[1:])

	switch opts.Subcommand {
	case "serve":
		if err := serve(opts.ConfigPath); err != nil {
			slog.Error("serve failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	case "migrate":
		if err := migrate(opts.ConfigPath); err != nil {
			slog.Error("migrate failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	case "backfill":
		if err := backfill(opts); err != nil {
			slog.Error("backfill failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	case "version":
		printVersion()
	default:
		slog.Error("unknown subcommand", slog.String("cmd", opts.Subcommand))
		slog.Info(
			"usage: olhojogo [serve|migrate|backfill|version] [-config path] [-from YYYY-MM-DD] [-to YYYY-MM-DD]",
		)
		os.Exit(exitUsage)
	}
}
