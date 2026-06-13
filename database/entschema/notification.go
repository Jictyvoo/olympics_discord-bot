package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Notification records an attempt to send an Alert to a channel. ID is a
// UUIDv7 minted by the UUIDMixin.
type Notification struct{ ent.Schema }

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("alert_id"),
		field.String("channel_id").Optional(),
		field.String("message_id").Optional(),
		field.String("status").NotEmpty().Default("pending"),
		field.String("checksum").MaxLen(64).Optional(),
		field.Time("sent_at").Optional().Default(time.Now),
	}
}

func (Notification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("alert_id"),
		index.Fields("checksum"),
	}
}

func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("alert", Alert.Type).Ref("notifications").Field("alert_id").Unique().Required(),
	}
}

func (Notification) Mixin() []ent.Mixin {
	return []ent.Mixin{customixins.UUIDMixin{}, customixins.TimestampsMixin{}}
}
