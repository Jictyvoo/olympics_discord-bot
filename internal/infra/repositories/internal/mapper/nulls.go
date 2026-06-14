package mapper

import (
	"time"

	"github.com/jictyvoo/olhojogo/pkg/idgen"
)

func NullStr(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func NullStrPtr(v any) *string {
	s := NullStr(v)
	if s == "" {
		return nil
	}
	return new(s)
}

func NullInt(v any) *int {
	if v == nil {
		return nil
	}
	switch n := v.(type) {
	case int64:
		i := int(n)
		return new(i)
	case int:
		return new(n)
	}
	return nil
}

func NullTime(v any) *time.Time {
	if v == nil {
		return nil
	}
	if t, ok := v.(time.Time); ok {
		return new(t)
	}
	return nil
}

func TimeOrZero(v any) time.Time {
	if t := NullTime(v); t != nil {
		return t.UTC()
	}
	return time.Time{}
}

// IDFromBytes returns Zero on nil, non-byte input, or a slice under 16 bytes.
func IDFromBytes(v any) idgen.CanonicalID {
	switch b := v.(type) {
	case []byte:
		return idgen.FromBytes(b)
	case idgen.CanonicalID:
		return b
	}
	return idgen.Zero
}

func IDToBytes(id idgen.CanonicalID) []byte {
	return id.Bytes()
}

// NullableID returns nil for a nil pointer, or the ID bytes otherwise.
func NullableID(id *idgen.CanonicalID) any {
	if id == nil {
		return nil
	}
	return id.Bytes()
}

// IDPtrFromBytes returns a *CanonicalID for a non-empty []byte column, or nil
// when the column is NULL.
func IDPtrFromBytes(v any) *idgen.CanonicalID {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok || len(b) < 16 {
		return nil
	}
	id := idgen.FromBytes(b)
	return new(id)
}
