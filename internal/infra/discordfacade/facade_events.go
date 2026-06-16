package discordfacade

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// External events require a non-empty location and an end strictly after the
// start; these cover fixtures whose venue or end time is unknown.
const (
	defaultEventLocation = "TBD"
	defaultEventDuration = 2 * time.Hour
)

func (c *Client) CreateScheduledEvent(
	guildID string,
	in ScheduledEventInput,
) (string, error) {
	params := buildParams(in, discordgo.GuildScheduledEventStatusScheduled)
	ev, err := c.session.GuildScheduledEventCreate(guildID, params)
	if err != nil {
		return "", fmt.Errorf("discordfacade: create event: %w", err)
	}
	return ev.ID, nil
}

func (c *Client) UpdateScheduledEvent(
	guildID, eventID string,
	in ScheduledEventInput,
) error {
	params := buildParams(in, 0) // preserve existing status
	_, err := c.session.GuildScheduledEventEdit(guildID, eventID, params)
	if err != nil {
		return fmt.Errorf("discordfacade: update event: %w", err)
	}
	return nil
}

func (c *Client) CancelScheduledEvent(guildID, eventID string) error {
	params := &discordgo.GuildScheduledEventParams{
		Status: discordgo.GuildScheduledEventStatusCanceled,
	}
	_, err := c.session.GuildScheduledEventEdit(guildID, eventID, params)
	if err != nil {
		return fmt.Errorf("discordfacade: cancel event: %w", err)
	}
	return nil
}

func buildParams(
	in ScheduledEventInput,
	status discordgo.GuildScheduledEventStatus,
) *discordgo.GuildScheduledEventParams {
	start := in.StartsAt
	end := in.EndsAt
	if !end.After(start) {
		end = start.Add(defaultEventDuration)
	}

	p := &discordgo.GuildScheduledEventParams{
		Name:               in.Name,
		Description:        in.Description,
		ScheduledStartTime: &start,
		ScheduledEndTime:   &end,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		Status:             status,
	}
	if in.ChannelID != "" {
		p.ChannelID = in.ChannelID
		p.EntityType = discordgo.GuildScheduledEventEntityTypeVoice
		return p
	}

	// External events require entity_metadata.location to be present.
	location := in.Location
	if location == "" {
		location = defaultEventLocation
	}
	p.EntityType = discordgo.GuildScheduledEventEntityTypeExternal
	p.EntityMetadata = &discordgo.GuildScheduledEventEntityMetadata{Location: location}
	return p
}
