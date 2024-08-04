package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	"github.com/jictyvoo/olympics_data_fetcher/build/entschema/customixins"
)

// Results holds the schema definition for the Results entity.
type Results struct {
	ent.Schema
}

// Fields of the Results.
func (Results) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Unique(),
		field.Uint64("competitor_id").Immutable(),
		field.Uint64("event_id").Immutable(),
		field.String("position").Optional(),
		field.String("mark").Optional(),
		field.String("medal_type").Optional(),
		field.String("irm"),
	}
}

// Edges of the Results.
func (Results) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("competitor", Competitors.Type).
			Ref("results").Field("competitor_id").
			Unique().Required().Immutable(),
		edge.From("olympic_events", OlympicEvent.Type).
			Ref("results").Field("event_id").
			Unique().Required().Immutable(),
	}
}

// Indexes of the Results.
func (Results) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("competitor_id", "event_id").
			Unique(),
	}
}

// Mixin of the Results.
func (Results) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.TimestampsMixin{},
	}
}
