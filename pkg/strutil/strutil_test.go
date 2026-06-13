package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAlfaNum(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "Identical strings",
			a:        "abc123",
			b:        "abc123",
			expected: true,
		},
		{
			name:     "Different strings with same alphanumeric characters",
			a:        "abc123",
			b:        "a.b-c 123!",
			expected: true,
		},
		{
			name:     "Different strings with different alphanumeric characters",
			a:        "abc123",
			b:        "def456",
			expected: false,
		},
		{
			name:     "One string is a prefix of the other",
			a:        "abc123",
			b:        "abc",
			expected: false,
		},
		{
			name:     "Strings with different cases",
			a:        "abc123",
			b:        "ABC123",
			expected: true,
		},
		{
			name:     "Strings with non-alphanumeric characters only",
			a:        "!@#$%^&*()",
			b:        "()&*^%$#@!",
			expected: true,
		},
		{
			name:     "Empty strings",
			a:        "",
			b:        "",
			expected: true,
		},
		{
			name:     "One empty string and one non-empty string",
			a:        "",
			b:        "abc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := EqualAlfaNum(tt.a, tt.b)
				assert.Equal(t, tt.expected, result)
			},
		)
	}
}
