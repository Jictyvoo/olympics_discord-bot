package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// OlympicDiscipline holds the schema definition for the OlympicDiscipline (OlympicDiscipline) entity.
type OlympicDiscipline struct {
	ent.Schema
}

// Fields of the OlympicDiscipline.
func (OlympicDiscipline) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.String("name").Unique(),
		field.String("description").Optional(),
	}
}

// Indexes of the OlympicDiscipline.
func (OlympicDiscipline) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("name").Unique(),
	}
}

// Edges of the OlympicDiscipline.
func (OlympicDiscipline) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("olympic_events", OlympicEvent.Type),
	}
}
