package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Result holds the outcome for a single participant within a Fixture.
type Result struct{ ent.Schema }

func (Result) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("fixture_id"),
		uuidfield.Field("participant_id"),
		field.Int("position").Optional().Nillable(),
		field.String("score").Optional(),
		field.String("raw_mark").Optional(),
		field.String("outcome").Optional(),
	}
}

func (Result) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("fixture_id", "participant_id").Unique(),
		index.Fields("participant_id"),
	}
}

func (Result) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fixture", Fixture.Type).Ref("results").Field("fixture_id").Unique().Required(),
	}
}

func (Result) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
