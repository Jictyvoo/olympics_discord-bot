package entities

import (
	"testing"
)

func TestGetLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected Language
	}{
		{"en", languageList[0]},
		{"ENG", languageList[0]},
		{"fr", languageList[1]},
		{"FRA", languageList[1]},
		{"de", languageList[2]},
		{"DEU", languageList[2]},
		{"it", languageList[3]},
		{"ITA", languageList[3]},
		{"pt", languageList[4]},
		{"POR", languageList[4]},
		{"es", languageList[5]},
		{"SPA", languageList[5]},
		{"ja", languageList[6]},
		{"JPN", languageList[6]},
		{"ar", languageList[7]},
		{"zh", languageList[8]},
		{"CHI", languageList[8]},
		{"hi", languageList[9]},
		{"ko", languageList[10]},
		{"KOR", languageList[10]},
		{"ru", languageList[11]},
		{"RUS", languageList[11]},
		{"unknown", languageList[0]},
	}

	for _, test := range tests {
		result := GetLanguage(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected %+v but got %+v", test.input, test.expected, result)
		}
	}
}
