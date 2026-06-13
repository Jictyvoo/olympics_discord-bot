package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Country holds the schema definition for the Country entity.
// Seeded once from countries.json; not owned by any provider.
type Country struct {
	ent.Schema
}

func (Country) Fields() []ent.Field {
	return []ent.Field{
		field.String("iso2").MaxLen(2).Unique().NotEmpty(),
		field.String("iso3").MaxLen(3).Unique().NotEmpty(),
		field.String("ioc_code").MaxLen(3).Optional(),
		field.String("name").NotEmpty(),
		field.Int("code_num").Optional(),
		field.Int64("population").Optional(),
		field.Float("area_km2").Optional(),
		field.Float("gdp_usd").Optional(),
	}
}

func (Country) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ioc_code"),
		index.Fields("name"),
	}
}

func (Country) Edges() []ent.Edge { return nil }
