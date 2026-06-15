package discordfacade

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ScheduledEventInput carries the data needed to create or update a Discord scheduled event.
type ScheduledEventInput struct {
	Name        string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	ChannelID   string
	Location    string
}

// Client wraps bwmarrin/discordgo for the operations we actually use.
// Callers declare their own narrower interfaces over these methods.
type Client struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *Client {
	return &Client{session: session}
}

func (c *Client) ResolveChannel(_ context.Context, guildID, channelName string) (string, error) {
	channels, err := c.session.GuildChannels(guildID)
	if err != nil {
		return "", fmt.Errorf("discordfacade: list channels: %w", err)
	}
	for _, ch := range channels {
		if strings.EqualFold(ch.Name, channelName) {
			return ch.ID, nil
		}
	}
	ch, err := c.session.GuildChannelCreate(guildID, channelName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return "", fmt.Errorf("discordfacade: create channel: %w", err)
	}
	return ch.ID, nil
}

func (c *Client) Send(channelID, content string) (string, error) {
	msg, err := c.session.ChannelMessageSend(channelID, content)
	if err != nil {
		return "", fmt.Errorf("discordfacade: send message: %w", err)
	}
	return msg.ID, nil
}

// InstallRouter wires the router's dispatcher into the session and registers its
// guild-scoped application command, keeping all discordgo calls in this package.
func (c *Client) InstallRouter(r *Router, guildID string) error {
	c.session.AddHandler(r.Handle)
	appID := c.session.State.User.ID
	if _, err := c.session.ApplicationCommandCreate(
		appID,
		guildID,
		r.ApplicationCommand(),
	); err != nil {
		return fmt.Errorf("discordfacade: register commands: %w", err)
	}
	return nil
}
