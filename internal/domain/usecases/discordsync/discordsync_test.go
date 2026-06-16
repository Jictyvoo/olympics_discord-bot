package discordsync

import (
	"database/sql"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/discordfacade"
)

func mkFixture(checksum string, status eventcore.FixtureStatus) eventcore.Fixture {
	id := eventcore.NewID(eventcore.ProviderOlympics, "fx-"+checksum)
	return eventcore.Fixture{
		ID:       id,
		Name:     "Match " + checksum,
		StartsAt: time.Now().Add(time.Hour),
		EndsAt:   time.Now().Add(2 * time.Hour),
		Status:   status,
		Checksum: checksum,
	}
}

func TestDiscordSync_CreatesWhenMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("c1", eventcore.FixtureScheduled)

	reader := NewMockFixtureReader(ctrl)
	reader.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{}, sql.ErrNoRows)
	var upserted eventcore.DiscordEvent
	repo.EXPECT().
		UpsertDiscordEvent(gomock.Any()).
		Do(func(de eventcore.DiscordEvent) { upserted = de }).
		Return(nil)

	fac := NewMockScheduledEventFacade(ctrl)
	fac.EXPECT().ListScheduledEvents(gomock.Any()).Return(nil, nil)
	fac.EXPECT().
		CreateScheduledEvent(gomock.Any(), gomock.Any()).
		Return("discord-1", nil)

	ds := New(reader, repo, fac, "guild", 0)
	if err := ds.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if upserted.DiscordEventID != "discord-1" {
		t.Fatalf("expected upsert with discord-1; got %+v", upserted)
	}
}

func TestDiscordSync_AdoptsExistingEventInsteadOfCreating(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("dup", eventcore.FixtureScheduled)

	reader := NewMockFixtureReader(ctrl)
	reader.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{}, sql.ErrNoRows)
	var upserted eventcore.DiscordEvent
	repo.EXPECT().
		UpsertDiscordEvent(gomock.Any()).
		Do(func(de eventcore.DiscordEvent) { upserted = de }).
		Return(nil)

	// Discord already has an event for this fixture (same description) -> adopt it.
	existing := discordfacade.ScheduledEvent{
		ID:          "evt-existing",
		Description: buildEventInput(f).Description,
	}
	fac := NewMockScheduledEventFacade(ctrl)
	fac.EXPECT().
		ListScheduledEvents(gomock.Any()).
		Return([]discordfacade.ScheduledEvent{existing}, nil)
	// CreateScheduledEvent must NOT be called: no EXPECT() => any call fails.

	ds := New(reader, repo, fac, "guild", 0)
	if err := ds.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if upserted.DiscordEventID != "evt-existing" {
		t.Fatalf("expected to adopt evt-existing; got %+v", upserted)
	}
}

func TestDiscordSync_UpdatesWhenChecksumChanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("c2", eventcore.FixtureScheduled)

	reader := NewMockFixtureReader(ctrl)
	reader.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{
			FixtureID: f.ID, DiscordEventID: "evt-99", LastChecksum: "stale",
		}, nil)
	var upserted eventcore.DiscordEvent
	repo.EXPECT().
		UpsertDiscordEvent(gomock.Any()).
		Do(func(de eventcore.DiscordEvent) { upserted = de }).
		Return(nil)

	fac := NewMockScheduledEventFacade(ctrl)
	var updatedEventID string
	fac.EXPECT().
		UpdateScheduledEvent(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_, eventID string, _ discordfacade.ScheduledEventInput) {
			updatedEventID = eventID
		}).
		Return(nil)

	ds := New(reader, repo, fac, "guild", 0)
	if err := ds.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if updatedEventID != "evt-99" {
		t.Fatalf("expected update on evt-99; got id=%q", updatedEventID)
	}
	if upserted.LastChecksum != "c2" {
		t.Fatalf("expected upsert with new checksum; got %+v", upserted)
	}
}

func TestDiscordSync_SkipsWhenChecksumUnchanged(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("c3", eventcore.FixtureScheduled)

	reader := NewMockFixtureReader(ctrl)
	reader.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{
			FixtureID: f.ID, DiscordEventID: "evt-7", LastChecksum: "c3",
		}, nil)

	// No facade calls expected: NewMockScheduledEventFacade with no EXPECT()
	// will fail the test if any method is invoked.
	fac := NewMockScheduledEventFacade(ctrl)

	ds := New(reader, repo, fac, "guild", 0)
	if err := ds.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func TestDiscordSync_SkipsPastFixture(t *testing.T) {
	ctrl := gomock.NewController(t)
	// A historical (already kicked-off) fixture from the backfill.
	f := mkFixture("past", eventcore.FixtureFinished)
	f.StartsAt = time.Now().Add(-3 * time.Hour)
	f.EndsAt = time.Now().Add(-time.Hour)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{}, sql.ErrNoRows)

	// No facade or upsert calls: Discord cannot schedule events in the past.
	fac := NewMockScheduledEventFacade(ctrl)

	ds := New(NewMockFixtureReader(ctrl), repo, fac, "guild", 0)
	ds.On(f)
}

func TestDiscordSync_CancelsWhenFixtureCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	f := mkFixture("c4", eventcore.FixtureCancelled)

	reader := NewMockFixtureReader(ctrl)
	reader.EXPECT().
		ListFixturesStartingBefore(gomock.Any()).
		Return([]eventcore.Fixture{f}, nil)

	repo := NewMockDiscordEventRepo(ctrl)
	repo.EXPECT().
		GetDiscordEventByFixture(f.ID, gomock.Any()).
		Return(eventcore.DiscordEvent{
			FixtureID: f.ID, DiscordEventID: "evt-cancel-me", LastChecksum: "c4",
		}, nil)
	var gotStatus eventcore.DiscordEventStatus
	repo.EXPECT().
		UpdateDiscordEventStatus(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ eventcore.CanonicalID, _ string, st eventcore.DiscordEventStatus) {
			gotStatus = st
		}).
		Return(nil)

	fac := NewMockScheduledEventFacade(ctrl)
	fac.EXPECT().
		CancelScheduledEvent(gomock.Any(), gomock.Any()).
		Return(nil)

	ds := New(reader, repo, fac, "guild", 0)
	if err := ds.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if gotStatus != eventcore.DiscordEventCancelled {
		t.Fatalf("expected cancelled status update; got %v", gotStatus)
	}
}

func TestDiscordSync_DefaultHorizon(t *testing.T) {
	ctrl := gomock.NewController(t)
	ds := New(
		NewMockFixtureReader(ctrl),
		NewMockDiscordEventRepo(ctrl),
		NewMockScheduledEventFacade(ctrl),
		"g", 0,
	)
	if ds.horizon == 0 {
		t.Fatal("default horizon must be non-zero")
	}
}
