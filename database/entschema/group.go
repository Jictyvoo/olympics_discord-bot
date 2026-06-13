package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Group is an optional grouping within a Stage.
type Group struct{ ent.Schema }

func (Group) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("name").NotEmpty(),
		uuidfield.Field("stage_id"),
	}
}

func (Group) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "external_key").Unique(),
		index.Fields("stage_id"),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stage", Stage.Type).
			Ref("groups").Field("stage_id").Unique().Required(),
	}
}

func (Group) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
