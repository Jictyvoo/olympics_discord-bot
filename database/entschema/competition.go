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

// Competition is the top-level tournament concept. ID = idgen.From(provider, key).
type Competition struct{ ent.Schema }

func (Competition) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("code").Optional(),
		field.String("name").NotEmpty(),
		field.String("discipline").Optional(),
	}
}

func (Competition) Indexes() []ent.Index {
	return []ent.Index{index.Fields("provider_id", "external_key").Unique()}
}

func (Competition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("seasons", Season.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Competition) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
