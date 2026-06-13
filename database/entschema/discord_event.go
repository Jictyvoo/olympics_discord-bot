package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// DiscordEvent tracks a Discord Scheduled Event created for a Fixture.
type DiscordEvent struct{ ent.Schema }

func (DiscordEvent) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("fixture_id"),
		field.String("guild_id").NotEmpty(),
		field.String("discord_event_id").Unique().NotEmpty(),
		field.String("status").NotEmpty().Default("scheduled"),
		field.String("last_checksum").MaxLen(64).Optional(),
	}
}

func (DiscordEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("fixture_id", "guild_id").Unique(),
		index.Fields("discord_event_id").Unique(),
	}
}

func (DiscordEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fixture", Fixture.Type).
			Ref("discord_events").
			Field("fixture_id").
			Unique().
			Required(),
	}
}

func (DiscordEvent) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
