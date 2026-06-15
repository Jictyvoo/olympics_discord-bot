package discordfacade

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// CommandInvocation decouples consumers from discordgo's interaction types.
type CommandInvocation struct {
	Name      string
	Sub       string
	GuildID   string
	ChannelID string
	UserID    string
	Options   map[string]string
}

type CommandHandler interface {
	Handle(ctx context.Context, inv CommandInvocation) (reply string, err error)
}

func (c *Client) RegisterCommands(handler CommandHandler, guildID string) error {
	c.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		inv := parseInvocation(i)
		reply, err := handler.Handle(context.Background(), inv)
		if err != nil {
			reply = err.Error()
		}
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: reply,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	})

	appID := c.session.State.User.ID
	if _, err := c.session.ApplicationCommandCreate(appID, guildID, notifyCommand()); err != nil {
		return fmt.Errorf("discordfacade: register commands: %w", err)
	}
	return nil
}

func parseInvocation(i *discordgo.InteractionCreate) CommandInvocation {
	data := i.ApplicationCommandData()
	inv := CommandInvocation{
		Name:      data.Name,
		GuildID:   i.GuildID,
		ChannelID: i.ChannelID,
		Options:   map[string]string{},
	}
	if i.Member != nil && i.Member.User != nil {
		inv.UserID = i.Member.User.ID
	} else if i.User != nil {
		inv.UserID = i.User.ID
	}
	for _, opt := range data.Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			inv.Sub = opt.Name
			for _, subOpt := range opt.Options {
				inv.Options[subOpt.Name] = subOpt.StringValue()
			}
		}
	}
	return inv
}

func notifyCommand() *discordgo.ApplicationCommand {
	kindOption := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "kind",
		Description: "Subscription kind",
		Required:    true,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{Name: "all_results", Value: "all_results"},
			{Name: "country", Value: "country"},
			{Name: "discipline", Value: "discipline"},
		},
	}
	valueOption := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "value",
		Description: "Country or discipline code",
		Required:    false,
	}
	return &discordgo.ApplicationCommand{
		Name:        "notify",
		Description: "Manage result notification subscriptions",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add a subscription",
				Options:     []*discordgo.ApplicationCommandOption{kindOption, valueOption},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a subscription",
				Options:     []*discordgo.ApplicationCommandOption{kindOption, valueOption},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List subscriptions for this server",
			},
		},
	}
}
