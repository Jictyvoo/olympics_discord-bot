package discordfacade

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestBuildParams_ExternalDefaultsLocationAndEnd(t *testing.T) {
	start := time.Date(2026, 6, 18, 18, 0, 0, 0, time.UTC)
	// EndsAt == StartsAt and no location: both must be made valid for Discord.
	in := ScheduledEventInput{Name: "Brazil vs Argentina", StartsAt: start, EndsAt: start}

	p := buildParams(in, discordgo.GuildScheduledEventStatusScheduled)

	if p.EntityType != discordgo.GuildScheduledEventEntityTypeExternal {
		t.Fatalf("expected external entity type, got %v", p.EntityType)
	}
	if p.EntityMetadata == nil || p.EntityMetadata.Location != defaultEventLocation {
		t.Fatalf("external event must carry a non-empty location, got %+v", p.EntityMetadata)
	}
	if !p.ScheduledEndTime.After(*p.ScheduledStartTime) {
		t.Fatalf("end %v must be after start %v", p.ScheduledEndTime, p.ScheduledStartTime)
	}
}

func TestBuildParams_ExternalKeepsGivenLocation(t *testing.T) {
	start := time.Date(2026, 6, 18, 18, 0, 0, 0, time.UTC)
	in := ScheduledEventInput{
		Name:     "Final",
		StartsAt: start,
		EndsAt:   start.Add(time.Hour),
		Location: "Maracanã",
	}

	p := buildParams(in, discordgo.GuildScheduledEventStatusScheduled)

	if p.EntityMetadata.Location != "Maracanã" {
		t.Errorf("location = %q, want Maracanã", p.EntityMetadata.Location)
	}
	if !p.ScheduledEndTime.Equal(start.Add(time.Hour)) {
		t.Errorf("a valid end time must be preserved, got %v", p.ScheduledEndTime)
	}
}

func TestBuildParams_VoiceUsesChannel(t *testing.T) {
	start := time.Date(2026, 6, 18, 18, 0, 0, 0, time.UTC)
	in := ScheduledEventInput{
		Name:      "Watch party",
		StartsAt:  start,
		EndsAt:    start,
		ChannelID: "chan-1",
	}

	p := buildParams(in, discordgo.GuildScheduledEventStatusScheduled)

	if p.EntityType != discordgo.GuildScheduledEventEntityTypeVoice || p.ChannelID != "chan-1" {
		t.Fatalf("expected voice event on channel, got type=%v chan=%q", p.EntityType, p.ChannelID)
	}
	if p.EntityMetadata != nil {
		t.Errorf("voice events must not set entity metadata, got %+v", p.EntityMetadata)
	}
}
