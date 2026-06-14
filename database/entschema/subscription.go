package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
)

// Subscription is a Discord user's mention request scoped to a guild.
type Subscription struct{ ent.Schema }

func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("guild_id").MaxLen(32).NotEmpty(),
		field.String("user_id").MaxLen(32).NotEmpty(),
		field.String("kind").MaxLen(32).NotEmpty(),
		field.String("value").MaxLen(64).Optional(),
	}
}

func (Subscription) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("guild_id", "user_id", "kind", "value").Unique(),
		index.Fields("guild_id"),
	}
}

func (Subscription) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
