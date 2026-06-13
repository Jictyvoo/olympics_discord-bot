package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Standing tracks a Participant's rank within a Stage (league table).
type Standing struct{ ent.Schema }

func (Standing) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("stage_id"),
		uuidfield.Field("participant_id"),
		field.Int("rank").Default(0),
		field.Int("points").Default(0),
	}
}

func (Standing) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("stage_id", "participant_id").Unique(),
		index.Fields("stage_id", "rank"),
	}
}

func (Standing) Edges() []ent.Edge { return nil }

func (Standing) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
