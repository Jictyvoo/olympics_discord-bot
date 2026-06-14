package eventcore

import "testing"

func TestCountry_EmojiFlag(t *testing.T) {
	testCases := []struct {
		name string
		iso2 string
		want string
	}{
		{name: "lowercase iso2", iso2: "br", want: flagBR},
		{name: "uppercase iso2", iso2: "BR", want: flagBR},
		{name: "mixed case iso2", iso2: "Us", want: ":flag_us:"},
		{name: "empty iso2", iso2: "", want: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := Country{ISO2: tc.iso2}
			if got := c.EmojiFlag(); got != tc.want {
				t.Errorf("EmojiFlag() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCountry_IsThis(t *testing.T) {
	country := Country{
		ISO2:    "BR",
		ISO3:    bra,
		IOCCode: bra,
		Name:    brazil,
	}

	testCases := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "matches name", value: brazil, want: true},
		{name: "matches name case-insensitive", value: "brazil", want: true},
		{name: "matches iso2", value: "br", want: true},
		{name: "matches iso3", value: "bra", want: true},
		{name: "matches ioc code", value: bra, want: true},
		{name: "no match", value: "Argentina", want: false},
		{name: "empty value", value: "", want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := country.IsThis(tc.value); got != tc.want {
				t.Errorf("IsThis(%q) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

func TestCountry_IsThis_EmptyFieldsDoNotMatchEmptyValue(t *testing.T) {
	// A country with no ISO2 must not report an empty query as a match, since
	// every empty field would otherwise collapse to the same "" case.
	country := Country{Name: "Nowhere"}
	if country.IsThis("") {
		t.Error("IsThis(\"\") = true, want false for unset codes")
	}
}
