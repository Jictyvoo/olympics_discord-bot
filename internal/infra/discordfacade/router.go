package discordfacade

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// Router dispatches Discord interactions to subcommands, like http.ServeMux.
type Router struct {
	name        string
	description string
	subs        map[string]SubCommand
	logger      *slog.Logger
}

func NewRouter(name, description string, logger *slog.Logger) *Router {
	if logger == nil {
		logger = slog.Default()
	}
	return &Router{
		name:        name,
		description: description,
		subs:        make(map[string]SubCommand),
		logger:      logger,
	}
}

func (r *Router) Add(cmd SubCommand) *Router {
	r.subs[cmd.Spec().Name] = cmd
	return r
}

func (r *Router) ApplicationCommand() *discordgo.ApplicationCommand {
	opts := make([]*discordgo.ApplicationCommandOption, 0, len(r.subs))
	for _, sub := range r.subs {
		opts = append(opts, sub.Spec())
	}
	return &discordgo.ApplicationCommand{
		Name:        r.name,
		Description: r.description,
		Options:     opts,
	}
}

// Handle is the InteractionCreate handler registered with discordgo.
func (r *Router) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != r.name {
		return
	}

	opts := i.ApplicationCommandData().Options
	if len(opts) == 0 {
		return
	}
	sub, ok := r.subs[opts[0].Name]
	if !ok {
		r.logger.Warn("unknown subcommand", "name", opts[0].Name)
		return
	}

	ctx := context.Background()
	resp := interactionContext{session: s, interaction: i.Interaction}
	inv := invocationFrom(i)
	subOpts := newOptionMap(opts[0].Options)
	if err := sub.Handle(ctx, inv, subOpts, resp); err != nil {
		r.logger.Error("subcommand error", "cmd", opts[0].Name, "err", err)
	}
}

func invocationFrom(i *discordgo.InteractionCreate) Invocation {
	inv := Invocation{GuildID: i.GuildID}
	if i.Member != nil && i.Member.User != nil {
		inv.UserID = i.Member.User.ID
	} else if i.User != nil {
		inv.UserID = i.User.ID
	}
	return inv
}
