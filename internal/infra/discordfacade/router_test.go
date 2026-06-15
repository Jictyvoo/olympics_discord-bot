package discordfacade

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// stubCommands records the arguments of the last HandleCommand call.
type stubCommands struct {
	reply  string
	err    error
	action string
	guild  string
	user   string
	kind   string
	value  string
}

func (s *stubCommands) HandleCommand(
	action, guildID, userID, kind, value string,
) (string, error) {
	s.action, s.guild, s.user, s.kind, s.value = action, guildID, userID, kind, value
	return s.reply, s.err
}

// recordingResponder captures what a SubCommand sent.
type recordingResponder struct {
	content   string
	ephemeral bool
}

func (r *recordingResponder) Send(_ context.Context, content string, ephemeral bool) error {
	r.content, r.ephemeral = content, ephemeral
	return nil
}

func (r *recordingResponder) Edit(_ context.Context, content string) error {
	r.content = content
	return nil
}

func TestOptionMapString(t *testing.T) {
	opts := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: optKind, Type: discordgo.ApplicationCommandOptionString, Value: kindCntry},
	}
	m := newOptionMap(opts)
	if got := m.String(optKind); got != kindCntry {
		t.Errorf("String(kind) = %q, want %q", got, kindCntry)
	}
	if got := m.String("missing"); got != "" {
		t.Errorf("String(missing) = %q, want empty", got)
	}
}

func TestRouterApplicationCommand(t *testing.T) {
	svc := &stubCommands{}
	r := NewRouter("notify", "manage subs", nil).
		Add(AddCmd(svc)).
		Add(RemoveCmd(svc)).
		Add(ListCmd(svc))

	cmd := r.ApplicationCommand()
	if cmd.Name != "notify" {
		t.Errorf("Name = %q, want notify", cmd.Name)
	}
	if len(cmd.Options) != 3 {
		t.Fatalf("len(Options) = %d, want 3", len(cmd.Options))
	}
	for _, opt := range cmd.Options {
		if opt.Type != discordgo.ApplicationCommandOptionSubCommand {
			t.Errorf("option %q type = %v, want SubCommand", opt.Name, opt.Type)
		}
	}
}

func TestSubcommandSpecs(t *testing.T) {
	svc := &stubCommands{}
	tests := []struct {
		name     string
		cmd      SubCommand
		wantName string
		wantOpts int
	}{
		{actionAdd, AddCmd(svc), actionAdd, 2},
		{actionRemove, RemoveCmd(svc), actionRemove, 2},
		{actionList, ListCmd(svc), actionList, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := tt.cmd.Spec()
			if spec.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", spec.Name, tt.wantName)
			}
			if spec.Type != discordgo.ApplicationCommandOptionSubCommand {
				t.Errorf("Type = %v, want SubCommand", spec.Type)
			}
			if len(spec.Options) != tt.wantOpts {
				t.Errorf("len(Options) = %d, want %d", len(spec.Options), tt.wantOpts)
			}
		})
	}
}

func TestAddCmdHandleForwardsAndReplies(t *testing.T) {
	svc := &stubCommands{reply: "subscribed"}
	resp := &recordingResponder{}
	inv := Invocation{GuildID: "g1", UserID: "u1"}
	opts := newOptionMap([]*discordgo.ApplicationCommandInteractionDataOption{
		{Name: optKind, Type: discordgo.ApplicationCommandOptionString, Value: kindCntry},
		{Name: optValue, Type: discordgo.ApplicationCommandOptionString, Value: "BRA"},
	})

	if err := AddCmd(svc).Handle(context.Background(), inv, opts, resp); err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if svc.action != actionAdd || svc.guild != "g1" || svc.user != "u1" ||
		svc.kind != kindCntry || svc.value != "BRA" {
		t.Errorf("forwarded args = %+v, want add/g1/u1/country/BRA", svc)
	}
	if resp.content != "subscribed" || !resp.ephemeral {
		t.Errorf("reply = %q ephemeral=%v, want subscribed/true", resp.content, resp.ephemeral)
	}
}
