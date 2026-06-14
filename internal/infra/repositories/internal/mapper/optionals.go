package mapper

import "github.com/jictyvoo/olhojogo/pkg/idgen"

// OptBytes returns the ID bytes for a non-nil pointer, nil otherwise. SQLite's
// generated params type uses any for nullable BLOBs.
func OptBytes(id *idgen.CanonicalID) any {
	if id == nil {
		return nil
	}
	return id.Bytes()
}

// OptString returns the string, or nil when empty, for nullable text columns
// whose SQLite param type is any.
func OptString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// FloatOrZero reads a nullable SQLite real column into a float64.
func FloatOrZero(v any) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}
