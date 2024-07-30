package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NotifiedEvent holds the schema definition for the NotifiedEvent entity.
type NotifiedEvent struct {
	ent.Schema
}

// Fields of the NotifiedEvent.
func (NotifiedEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.Uint64("event_id").Immutable(),
		field.String("event_sha256").MaxLen(255).Unique().NotEmpty(),
		field.String("status").NotEmpty(),
		field.Time("notified_at").Default(time.Now).Optional(),
	}
}

// Edges of the NotifiedEvent.
func (NotifiedEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("olympic_events", OlympicEvent.Type).
			Ref("notified_events").Field("event_id").
			Unique().Immutable().Required(),
	}
}
