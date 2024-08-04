package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olympics_data_fetcher/build/entschema/customixins"
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
		field.Bool("has_medal").Default(false),
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
		edge.To("results", Results.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("notified_events", NotifiedEvent.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("olympic_disciplines", OlympicDiscipline.Type).
			Ref("olympic_events").
			Field("discipline_id").
			Unique().
			Required(),
	}
}

// Mixin of the OlympicEvent.
func (OlympicEvent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.TimestampsMixin{},
	}
}
