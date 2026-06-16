package main

import (
	"flag"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

type syncOptions struct {
	Provider eventcore.ProviderID
	Date     time.Time
}

func parseSyncArgs(args []string) (syncOptions, error) {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)

	var (
		provider = fs.String("provider", eventcore.ProviderOlympics, "Provider code")
		date     = fs.String("date", "", "Day to fetch in YYYY-MM-DD (default: today UTC)")
	)
	if err := fs.Parse(args); err != nil {
		return syncOptions{}, err
	}

	opts := syncOptions{
		Provider: *provider,
		Date:     time.Now().UTC(),
	}
	if *date != "" {
		parsed, err := time.Parse(time.DateOnly, *date)
		if err != nil {
			return syncOptions{}, err
		}
		opts.Date = parsed
	}
	return opts, nil
}
