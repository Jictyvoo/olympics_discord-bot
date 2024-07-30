package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OlympicEvent holds the schema definition for the OlympicEvent (OlympicEvent) entity.
type OlympicEvent struct {
	ent.Schema
}

// Fields of the OlympicEvent.
func (OlympicEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.String("event_name"),
		field.Uint64("discipline_id"),
		field.String("phase"),
		field.Uint8("gender"),
		field.String("session_code"),
		field.Time("start_at").Default(time.Now),
		field.Time("end_at").Default(time.Now),
		field.String("status"),
	}
}

// Indexes of the OlympicEvent.
func (OlympicEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_name", "discipline_id", "phase", "gender", "session_code").Unique(),
	}
}

// Edges of the OlympicEvent.
func (OlympicEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("results", Results.Type),
		edge.To("notified_events", NotifiedEvent.Type),
		edge.From("olympic_disciplines", OlympicDiscipline.Type).
			Ref("olympic_events").
			Field("discipline_id").
			Unique().
			Required(),
	}
}
