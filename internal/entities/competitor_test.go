package entities

import (
	"math/rand"
	"slices"
)

func mockOlympicCompetitors() (competitors [7]OlympicCompetitors) {
	// Predefined lists of names, codes, and country codes
	var (
		names = []string{
			"Excalibur", "Rhino", "Volt", "Mag",
			"Nova", "Frost", "Ember", "Citrine",
			"Voruna", "Jade", "Yareli", "Saryn",
		}
		codes = []string{
			"Grineer", "Corpus", "Infested", "Orokin",
			"Tenno", "Sentient", "Stalker", "Murmur",
		}
		countryCodes = []string{"JPN", "BRA", "USA", "ITA", "GER"}
	)

	for index := 0; index < 7; index++ {
		name := names[rand.Intn(len(names))]
		age := rand.Intn(40-18) + 18
		countryCode := countryCodes[rand.Intn(len(countryCodes))]
		country := GetCountryByCode(countryCode)
		codeIndex := rand.Intn(len(codes))
		code := codes[codeIndex]
		codes = slices.Delete(codes, codeIndex, codeIndex+1)

		competitors[index] = OlympicCompetitors{
			ID:          Identifier(index + 1),
			Name:        name,
			Code:        code,
			Age:         uint8(age),
			CountryCode: countryCode,
			Country:     country,
		}
	}

	return competitors
}
