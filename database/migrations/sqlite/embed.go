package sqlite

import "embed"

// FS holds all SQLite migration files embedded at build time.
//
//go:embed *.sql
var FS embed.FS
