package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Competitors holds the schema definition for the Competitors (OlympicCompetitors) entity.
type Competitors struct {
	ent.Schema
}

// Fields of the Competitors.
func (Competitors) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique(),
		field.String("name"),
		field.String("code"),
		field.Uint64("country_id"),
	}
}

// Indexes of the Competitors.
func (Competitors) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("country_id"),
		index.Fields("code"),
		index.Fields("name", "code").Unique(),
		index.Fields("country_id", "name", "code").Unique(),
	}
}

// Edges of the Competitors.
func (Competitors) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("country_info", CountryInfo.Type).
			Ref("competitors").
			Field("country_id").
			Unique().
			Required(),
		edge.To("results", Results.Type),
	}
}
