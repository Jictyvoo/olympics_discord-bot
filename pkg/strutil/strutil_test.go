package strutil

import "testing"

const sampleAlfaNum = "abc123"

type equalAlfaNumCase struct {
	name     string
	a        string
	b        string
	expected bool
}

func runEqualAlfaNumCases(t *testing.T, tests []equalAlfaNumCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualAlfaNum(tt.a, tt.b); got != tt.expected {
				t.Fatalf("EqualAlfaNum(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestEqualAlfaNum(t *testing.T) {
	tests := []equalAlfaNumCase{
		{
			name:     "Identical strings",
			a:        sampleAlfaNum,
			b:        sampleAlfaNum,
			expected: true,
		},
		{
			name:     "Different strings with same alphanumeric characters",
			a:        sampleAlfaNum,
			b:        "a.b-c 123!",
			expected: true,
		},
		{
			name:     "Different strings with different alphanumeric characters",
			a:        sampleAlfaNum,
			b:        "def456",
			expected: false,
		},
		{
			name:     "One string is a prefix of the other",
			a:        sampleAlfaNum,
			b:        "abc",
			expected: false,
		},
		{
			name:     "Strings with different cases",
			a:        sampleAlfaNum,
			b:        "ABC123",
			expected: true,
		},
	}

	runEqualAlfaNumCases(t, tests)
}

func TestEqualAlfaNum_EdgeCases(t *testing.T) {
	tests := []equalAlfaNumCase{
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

	runEqualAlfaNumCases(t, tests)
}
