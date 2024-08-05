package entities

import "testing"

func TestEmojiFlag(t *testing.T) {
	tests := []struct {
		country CountryInfo
		want    string
	}{
		{CountryInfo{ISOCode: [2]string{"BR", "BRA"}}, "flag_br"},
		{CountryInfo{ISOCode: [2]string{"JP", "JPN"}}, "flag_jp"},
	}

	for _, tt := range tests {
		t.Run(
			tt.country.ISOCode[0], func(t *testing.T) {
				if got := tt.country.EmojiFlag(); got != tt.want {
					t.Errorf("EmojiFlag() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestIsThis(t *testing.T) {
	tests := []struct {
		country CountryInfo
		value   string
		want    bool
	}{
		{countriesData["BRA"], "Brazil", true},
		{countriesData["BRA"], "bra", true},
		{countriesData["BRA"], "BR", true},
		{countriesData["BRA"], "JPN", false},
		{countriesData["JPN"], "Japan", true},
		{countriesData["JPN"], "jp", true},
		{countriesData["JPN"], "JPN", true},
		{countriesData["JPN"], "BRA", false},
	}

	for _, tt := range tests {
		t.Run(
			tt.value, func(t *testing.T) {
				if got := tt.country.IsThis(tt.value); got != tt.want {
					t.Errorf("IsThis() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGetCountryByCode(t *testing.T) {
	tests := []struct {
		countryCode string
		want        CountryInfo
	}{
		{"BRA", countriesData["BRA"]},
		{"BR", countriesData["BRA"]},
		{"JPN", countriesData["JPN"]},
		{"JP", countriesData["JPN"]},
		{"ARG", countriesData["ARG"]},
	}

	for _, tt := range tests {
		t.Run(
			tt.countryCode, func(t *testing.T) {
				if got := GetCountryByCode(tt.countryCode); got != tt.want {
					t.Errorf("GetCountryByCode() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
