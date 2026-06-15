package discordfacade

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// SubscriptionCommands is the narrow domain port the notify subcommands depend
// on. *subscriptions.Service satisfies it structurally, keeping the domain
// package free of any discordfacade dependency.
type SubscriptionCommands interface {
	HandleCommand(
		action, guildID, userID, kind, value string,
	) (string, error)
}

const (
	actionAdd    = "add"
	actionRemove = "remove"
	actionList   = "list"

	optKind    = "kind"
	optValue   = "value"
	kindResult = "all_results"
	kindCntry  = "country"
	kindDisc   = "discipline"
)

// kindOption / valueOption describe the option schema shared by add and remove.
func kindOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        optKind,
		Description: "Subscription kind",
		Required:    true,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{Name: kindResult, Value: kindResult},
			{Name: kindCntry, Value: kindCntry},
			{Name: kindDisc, Value: kindDisc},
		},
	}
}

func valueOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        optValue,
		Description: "Country or discipline code",
		Required:    false,
	}
}

// addCmd handles /notify add.
type addCmd struct{ cmds SubscriptionCommands }

func (c *addCmd) Spec() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        actionAdd,
		Description: "Add a subscription",
		Options:     []*discordgo.ApplicationCommandOption{kindOption(), valueOption()},
	}
}

func (c *addCmd) Handle(
	ctx context.Context, inv Invocation, opts OptionMap, resp Responder,
) error {
	reply, err := c.cmds.HandleCommand(
		actionAdd, inv.GuildID, inv.UserID, opts.String(optKind), opts.String(optValue),
	)
	if err != nil {
		reply = err.Error()
	}
	return resp.Send(ctx, reply, true)
}

// removeCmd handles /notify remove.
type removeCmd struct{ cmds SubscriptionCommands }

func (c *removeCmd) Spec() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        actionRemove,
		Description: "Remove a subscription",
		Options:     []*discordgo.ApplicationCommandOption{kindOption(), valueOption()},
	}
}

func (c *removeCmd) Handle(
	ctx context.Context, inv Invocation, opts OptionMap, resp Responder,
) error {
	reply, err := c.cmds.HandleCommand(
		actionRemove, inv.GuildID, inv.UserID, opts.String(optKind), opts.String(optValue),
	)
	if err != nil {
		reply = err.Error()
	}
	return resp.Send(ctx, reply, true)
}

// listCmd handles /notify list.
type listCmd struct{ cmds SubscriptionCommands }

func (c *listCmd) Spec() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        actionList,
		Description: "List subscriptions for this server",
	}
}

func (c *listCmd) Handle(
	ctx context.Context, inv Invocation, _ OptionMap, resp Responder,
) error {
	reply, err := c.cmds.HandleCommand(actionList, inv.GuildID, inv.UserID, "", "")
	if err != nil {
		reply = err.Error()
	}
	return resp.Send(ctx, reply, true)
}

// AddCmd returns a SubCommand for /notify add.
//
//nolint:ireturn // factory returning consumer interface by design
func AddCmd(cmds SubscriptionCommands) SubCommand { return &addCmd{cmds: cmds} }

// RemoveCmd returns a SubCommand for /notify remove.
//
//nolint:ireturn // factory returning consumer interface by design
func RemoveCmd(cmds SubscriptionCommands) SubCommand { return &removeCmd{cmds: cmds} }

// ListCmd returns a SubCommand for /notify list.
//
//nolint:ireturn // factory returning consumer interface by design
func ListCmd(cmds SubscriptionCommands) SubCommand { return &listCmd{cmds: cmds} }
