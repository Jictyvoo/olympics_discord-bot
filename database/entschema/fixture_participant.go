package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// FixtureParticipant links a Participant to a Fixture with an optional role.
type FixtureParticipant struct{ ent.Schema }

func (FixtureParticipant) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("fixture_id"),
		uuidfield.Field("participant_id"),
		field.String("role").Optional(),
	}
}

func (FixtureParticipant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("fixture_id", "participant_id").Unique(),
		index.Fields("participant_id"),
	}
}

func (FixtureParticipant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fixture", Fixture.Type).
			Ref("fixture_participants").
			Field("fixture_id").
			Unique().
			Required(),
	}
}
