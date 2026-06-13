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

// Season groups Stages within a Competition.
type Season struct{ ent.Schema }

func (Season) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("name").NotEmpty(),
		field.Time("starts_on"),
		field.Time("ends_on"),
		uuidfield.Field("competition_id"),
	}
}

func (Season) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "external_key").Unique(),
		index.Fields("competition_id"),
	}
}

func (Season) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("competition", Competition.Type).
			Ref("seasons").Field("competition_id").Unique().Required(),
		edge.To("stages", Stage.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Season) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
