package customixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"

	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// UUIDMixin provides a UUID primary key column (BINARY(16) / BLOB / uuid).
// IDs are minted in Go via pkg/idgen.NewV7 at insert time; no DB-side default,
// which keeps the schema portable across SQLite, MariaDB and stock MySQL 8.
type UUIDMixin struct{ mixin.Schema }

func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{uuidfield.ID(nil)}
}
