package discordfacade

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type OptionMap map[string]*discordgo.ApplicationCommandInteractionDataOption

func (m OptionMap) String(key string) string {
	if opt, ok := m[key]; ok {
		return opt.StringValue()
	}
	return ""
}

func newOptionMap(opts []*discordgo.ApplicationCommandInteractionDataOption) OptionMap {
	m := make(OptionMap, len(opts))
	for _, o := range opts {
		m[o.Name] = o
	}
	return m
}

type Invocation struct {
	GuildID string
	UserID  string
}

// SubCommand is a self-describing command unit, analogous to http.Handler.
type SubCommand interface {
	Spec() *discordgo.ApplicationCommandOption
	Handle(ctx context.Context, inv Invocation, opts OptionMap, resp Responder) error
}
