package idgen

import "testing"

func TestFrom_Deterministic(t *testing.T) {
	testCases := []struct {
		name     string
		provider ProviderID
		key      string
		want     string
	}{
		{"same inputs always match", "olympics", "EV-001", From("olympics", "EV-001").String()},
		{"different provider differs", "fifa", "EV-001", ""},
		{"empty key", "olympics", "", ""},
	}

	first := From("olympics", "EV-001")
	testCases[0].want = first.String()
	testCases[1].want = From("fifa", "EV-001").String()
	testCases[2].want = From("olympics", "").String()

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			got := From(tCase.provider, tCase.key)
			if got.String() != tCase.want {
				t.Fatalf("got %q want %q", got.String(), tCase.want)
			}
		})
	}
}

func TestFrom_UniquePerPair(t *testing.T) {
	testCases := []struct {
		name string
		a, b ExternalID
	}{
		{"different key", ExternalID{"p", "k1"}, ExternalID{"p", "k2"}},
		{"different provider", ExternalID{"p1", "k"}, ExternalID{"p2", "k"}},
		{"separator prevents collision", ExternalID{"ab", "c"}, ExternalID{"a", "bc"}},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if tCase.a.Canonical() == tCase.b.Canonical() {
				t.Fatalf("collision: %+v and %+v produced the same ID", tCase.a, tCase.b)
			}
		})
	}
}

func TestCanonicalID_IsZero(t *testing.T) {
	testCases := []struct {
		name string
		id   CanonicalID
		want bool
	}{
		{"zero value", Zero, true},
		{"non-zero", From("x", "y"), false},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if got := tCase.id.IsZero(); got != tCase.want {
				t.Fatalf("IsZero() = %v, want %v", got, tCase.want)
			}
		})
	}
}
