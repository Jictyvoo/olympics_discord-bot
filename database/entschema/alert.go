package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Alert represents a notable event on a Fixture. ID is a UUIDv7 minted at
// notification time (UUIDMixin handles generation).
type Alert struct{ ent.Schema }

func (Alert) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("fixture_id"),
		field.String("kind").NotEmpty(),
	}
}

func (Alert) Indexes() []ent.Index {
	return []ent.Index{index.Fields("fixture_id")}
}

func (Alert) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fixture", Fixture.Type).Ref("alerts").Field("fixture_id").Unique().Required(),
		edge.To("notifications", Notification.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Alert) Mixin() []ent.Mixin {
	return []ent.Mixin{customixins.UUIDMixin{}, customixins.TimestampsMixin{}}
}
