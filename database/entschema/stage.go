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

// Stage is an ordered phase within a Season.
type Stage struct{ ent.Schema }

func (Stage) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("name").NotEmpty(),
		field.Int("ord").Default(0),
		uuidfield.Field("season_id"),
	}
}

func (Stage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "external_key").Unique(),
		index.Fields("season_id"),
	}
}

func (Stage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("season", Season.Type).
			Ref("stages").Field("season_id").Unique().Required(),
		edge.To("fixtures", Fixture.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("groups", Group.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Stage) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
