package discordsync

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const (
	defaultHorizonDays = 14
	hoursPerDay        = 24
)

// DiscordSync creates, updates, and cancels Discord scheduled events to mirror fixture state.
type DiscordSync struct {
	fixtures FixtureReader
	repo     DiscordEventRepo
	discord  ScheduledEventFacade
	guildID  string
	horizon  time.Duration // fixtures within this window are managed
}

func New(
	fixtures FixtureReader,
	repo DiscordEventRepo,
	discord ScheduledEventFacade,
	guildID string,
	horizon time.Duration,
) *DiscordSync {
	if horizon == 0 {
		horizon = defaultHorizonDays * hoursPerDay * time.Hour
	}
	return &DiscordSync{
		fixtures: fixtures,
		repo:     repo,
		discord:  discord,
		guildID:  guildID,
		horizon:  horizon,
	}
}

// On satisfies services.Observer[eventcore.Fixture]: it reconciles the single
// fixture's Discord scheduled event whenever the syncer persists it.
func (ds *DiscordSync) On(f eventcore.Fixture) {
	if err := ds.reconcile(f); err != nil {
		slog.Error(
			"discordsync: on-fixture",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	}
}

// Run reconciles Discord scheduled events with stored fixtures within the horizon.
func (ds *DiscordSync) Run() error {
	deadline := time.Now().Add(ds.horizon)
	fixtures, err := ds.fixtures.ListFixturesStartingBefore(deadline)
	if err != nil {
		return err
	}

	var errs []error
	for _, f := range fixtures {
		if err := ds.reconcile(f); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (ds *DiscordSync) reconcile(f eventcore.Fixture) error {
	discordEvent, err := ds.repo.GetDiscordEventByFixture(f.ID, ds.guildID)
	notFound := errors.Is(err, sql.ErrNoRows)
	if err != nil && !notFound {
		return err
	}

	// Cancel if fixture is cancelled or postponed.
	if f.Status == eventcore.FixtureCancelled || f.Status == eventcore.FixturePostponed {
		if !notFound && discordEvent.DiscordEventID != "" {
			if cancelErr := ds.discord.CancelScheduledEvent(
				ds.guildID,
				discordEvent.DiscordEventID,
			); cancelErr != nil {
				slog.Error("discordsync: cancel", slog.String("err", cancelErr.Error()))
				return cancelErr
			}
			return ds.repo.UpdateDiscordEventStatus(
				f.ID,
				ds.guildID,
				eventcore.DiscordEventCancelled,
			)
		}
		return nil
	}

	input := buildEventInput(f)

	if notFound || discordEvent.DiscordEventID == "" {
		// Create new.
		evtID, createErr := ds.discord.CreateScheduledEvent(ds.guildID, input)
		if createErr != nil {
			return createErr
		}
		return ds.repo.UpsertDiscordEvent(eventcore.DiscordEvent{
			FixtureID:      f.ID,
			GuildID:        ds.guildID,
			DiscordEventID: evtID,
			Status:         eventcore.DiscordEventScheduled,
			LastChecksum:   f.Checksum,
		})
	}

	// Update only if checksum changed.
	if discordEvent.LastChecksum == f.Checksum {
		return nil
	}
	if updateErr := ds.discord.UpdateScheduledEvent(
		ds.guildID,
		discordEvent.DiscordEventID,
		input,
	); updateErr != nil {
		return updateErr
	}
	discordEvent.LastChecksum = f.Checksum
	return ds.repo.UpsertDiscordEvent(discordEvent)
}
