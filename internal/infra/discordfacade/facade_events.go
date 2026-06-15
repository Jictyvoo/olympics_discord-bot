package discordfacade

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
	p := &discordgo.GuildScheduledEventParams{
		Name:               in.Name,
		Description:        in.Description,
		ScheduledStartTime: &in.StartsAt,
		ScheduledEndTime:   &in.EndsAt,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		Status:             status,
	}
	if in.ChannelID != "" {
		p.ChannelID = in.ChannelID
		p.EntityType = discordgo.GuildScheduledEventEntityTypeVoice
	} else {
		p.EntityType = discordgo.GuildScheduledEventEntityTypeExternal
		if in.Location != "" {
			p.EntityMetadata = &discordgo.GuildScheduledEventEntityMetadata{Location: in.Location}
		}
	}
	return p
}
