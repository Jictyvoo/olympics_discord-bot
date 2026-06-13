package uuidfield

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"

	"github.com/google/uuid"
)

// Field returns a required UUID-shaped column (BINARY(16) / BLOB / uuid).
func Field(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		SchemaType(schemaType())
}

// Unique returns a required, unique UUID-shaped column.
func Unique(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		Unique().
		SchemaType(schemaType())
}

// Optional returns a nullable UUID-shaped column.
func Optional(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		Optional().
		Nillable().
		SchemaType(schemaType())
}

// ID returns a primary key UUID with a Default(uuid.New) for runtime-generated
// rows. Pass a non-nil defaultExpr to also emit a DB-side default (e.g.
// UUID_v7() on MariaDB); pass nil to let the Go side mint the value alone.
func ID(defaultExpr *entsql.Annotation) ent.Field {
	f := field.UUID("id", uuid.UUID{}).
		Unique().
		Immutable().
		Default(uuid.New).
		SchemaType(schemaType())
	if defaultExpr != nil {
		f = f.Annotations(defaultExpr)
	}
	return f
}

// DeterministicID returns a primary key UUID-shaped column WITHOUT a default.
// Use for entities whose ID is derived from upstream data via idgen.From, so
// the application supplies the ID at insert time and re-syncs are idempotent.
func DeterministicID() ent.Field {
	return field.UUID("id", uuid.UUID{}).
		Unique().
		Immutable().
		SchemaType(schemaType())
}

func schemaType() map[string]string {
	return map[string]string{
		dialect.MySQL:    "BINARY(16)",
		dialect.Postgres: "uuid",
		dialect.SQLite:   "blob",
	}
}
