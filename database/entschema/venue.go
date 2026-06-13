package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Venue is a physical location where Fixtures take place.
type Venue struct{ ent.Schema }

func (Venue) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("city").Optional(),
		field.String("country_iso").Optional(),
	}
}

func (Venue) Indexes() []ent.Index {
	return []ent.Index{index.Fields("provider_id", "external_key").Unique()}
}

func (Venue) Edges() []ent.Edge { return nil }

func (Venue) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
