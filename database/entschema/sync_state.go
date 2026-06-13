package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SyncState checkpoints the last-known sync position per provider + scope.
type SyncState struct {
	ent.Schema
}

func (SyncState) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_id").NotEmpty(),
		field.String("scope").NotEmpty(), // e.g. "daily", "season:abc123"
		field.String("cursor").Optional(),
		field.Time("last_synced_at").Optional().Default(time.Now),
		field.String("last_error").Optional(),
	}
}

func (SyncState) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "scope").Unique(),
	}
}

func (SyncState) Edges() []ent.Edge { return nil }
