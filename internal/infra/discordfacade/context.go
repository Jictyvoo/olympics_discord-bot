package discordfacade

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Responder interface {
	Send(ctx context.Context, content string, ephemeral bool) error
	Edit(ctx context.Context, content string) error
}

type interactionContext struct {
	session     *discordgo.Session
	interaction *discordgo.Interaction
}

// ephemeral replies are visible only to the invoking user.
func (c interactionContext) Send(_ context.Context, content string, ephemeral bool) error {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	return c.session.InteractionRespond(c.interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flags,
		},
	})
}

func (c interactionContext) Edit(_ context.Context, content string) error {
	_, err := c.session.InteractionResponseEdit(c.interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return err
}
