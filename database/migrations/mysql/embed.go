package mysql

import "embed"

// FS holds all MySQL migration files embedded at build time.
//
//go:embed *.sql
var FS embed.FS
