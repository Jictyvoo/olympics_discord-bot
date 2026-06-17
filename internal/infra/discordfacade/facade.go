package discordfacade

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ScheduledEventInput struct {
	Name        string
	Description string
	StartsAt    time.Time
	EndsAt      time.Time
	ChannelID   string
	Location    string
}

// Client wraps discordgo; callers declare their own narrower interfaces over it.
type Client struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *Client {
	return &Client{session: session}
}

// ResolveChannel creates the named text channel when it does not already exist.
func (c *Client) ResolveChannel(guildID, channelName string) (string, error) {
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

// Edit lets a fixture's notification evolve in place instead of reposting.
func (c *Client) Edit(channelID, messageID, content string) error {
	if _, err := c.session.ChannelMessageEdit(channelID, messageID, content); err != nil {
		return fmt.Errorf("discordfacade: edit message: %w", err)
	}
	return nil
}

// InstallRouter wires the dispatcher and registers its guild-scoped command,
// keeping all discordgo calls inside this package.
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
