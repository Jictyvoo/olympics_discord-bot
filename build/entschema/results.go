package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Results holds the schema definition for the Results entity.
type Results struct {
	ent.Schema
}

// Fields of the Results.
func (Results) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.Uint64("competitor_id"),
		field.Uint64("event_id"),
		field.String("position").Optional(),
		field.String("mark").Optional(),
		field.String("medal_type").Optional(),
		field.String("irm"),
	}
}

// Edges of the Results.
func (Results) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("competitor", Competitors.Type).Ref("results").Field("competitor_id").Unique().Required(),
		edge.From("olympic_events", OlympicEvent.Type).Ref("results").Field("event_id").Unique().Required(),
	}
}
