package customixins

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimestampsMixin struct {
	mixin.Schema
}

func (TimestampsMixin) Fields() []ent.Field {
	defaultDatetime := entsql.DefaultExprs(
		map[string]string{
			dialect.MySQL:    "NOW()",
			dialect.SQLite:   "DATETIME('now')",
			dialect.Postgres: "NOW()",
		},
	)
	return []ent.Field{
		field.Time("created_at").Immutable().Default(time.Now).Annotations(defaultDatetime),
		field.Time("updated_at").
			UpdateDefault(time.Now).
			Default(time.Now).
			Annotations(defaultDatetime),
	}
}
