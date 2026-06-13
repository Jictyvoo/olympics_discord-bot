package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Participant is a team or athlete participating in Fixtures.
type Participant struct{ ent.Schema }

func (Participant) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("kind").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("code").Optional(),
		field.String("country_iso").Optional(),
		field.String("gender").Optional(),
	}
}

func (Participant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "external_key").Unique(),
		index.Fields("country_iso"),
	}
}

func (Participant) Edges() []ent.Edge { return nil }

func (Participant) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
