package mapper

import (
	"database/sql"
	"time"

	"github.com/jictyvoo/olhojogo/pkg/idgen"
)

// NSStr returns a NullString that's valid iff s is non-empty.
func NSStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// NSInt returns a NullInt64 that's valid iff n != 0.
func NSInt(n int64) sql.NullInt64 {
	if n == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: n, Valid: true}
}

// NSFloat returns a NullFloat64 that's valid iff f != 0.
func NSFloat(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

// NSTime returns a NullTime that's valid iff t is non-zero.
func NSTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t.UTC(), Valid: true}
}

// NSIntFromPtr returns a NullInt64 from a *int (nil → invalid).
func NSIntFromPtr(p *int) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}

// NSStrFromID returns a NullString carrying the 16 raw ID bytes (sqlc treats
// MySQL's nullable BINARY columns as NullString).
func NSStrFromID(id *idgen.CanonicalID) sql.NullString {
	if id == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(id.Bytes()), Valid: true}
}

// IDFromNullStr decodes a NullString back into a *CanonicalID. Empty → nil.
func IDFromNullStr(ns sql.NullString) *idgen.CanonicalID {
	if !ns.Valid || len(ns.String) == 0 {
		return nil
	}
	id := idgen.FromBytes([]byte(ns.String))
	return new(id)
}

// IntFromNull64 returns an *int from a NullInt64 (invalid → nil).
func IntFromNull64(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return new(v)
}
