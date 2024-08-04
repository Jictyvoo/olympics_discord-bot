package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olympics_data_fetcher/build/entschema/customixins"
)

// CountryInfo holds the schema definition for the CountryInfo entity.
type CountryInfo struct {
	ent.Schema
}

// Fields of the CountryInfo.
func (CountryInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique().Immutable(),
		field.String("code"),
		field.String("name"),
		field.String("code_num"),
		field.String("iso_code_len2").Optional(),
		field.String("iso_code_len3"),
		field.String("ioc_code").Unique().MaxLen(3),
		field.Uint64("population").Optional(),
		field.Float("area_km2").Optional(),
		field.String("gdp_usd").Optional(),
	}
}

// Indexes of the CountryInfo.
func (CountryInfo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("ioc_code"),
	}
}

// Edges of the CountryInfo.
func (CountryInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("competitors", Competitors.Type).
			Annotations(entsql.OnDelete(entsql.Restrict)),
	}
}

// Mixin of the CountryInfo.
func (CountryInfo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.TimestampsMixin{},
	}
}
