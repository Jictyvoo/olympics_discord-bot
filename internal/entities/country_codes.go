package entities

import (
	_ "embed"
	"encoding/json"
	"strings"
)

var (
	//go:embed country_data.json
	countriesFileContents []byte
	countriesData         = map[string]CountryInfo{}
)

func init() {
	if err := json.Unmarshal(countriesFileContents, &countriesData); err != nil {
		panic(err)
	}
}

type CountryInfo struct {
	Name       string
	CodeNum    string
	ISOCode    [2]string
	IOCCode    string
	Population uint64
	AreaKm2    float64
	GDPUSD     string
}

func (c CountryInfo) EmojiFlag() string {
	return "flag_" + strings.ToLower(c.ISOCode[0])
}

func (c CountryInfo) IsThis(value string) bool {
	switch value = strings.ToLower(value); value {
	case strings.ToLower(c.Name),
		strings.ToLower(c.IOCCode),
		strings.ToLower(c.ISOCode[0]),
		strings.ToLower(c.ISOCode[1]):
		return true
	}
	return false
}

func GetCountryByCode(countryCode string) CountryInfo {
	found, ok := countriesData[countryCode]
	if ok {
		return found
	}

	for _, country := range countriesData {
		if country.IOCCode == countryCode || country.ISOCode[0] == countryCode ||
			country.ISOCode[1] == countryCode {
			return country
		}
	}

	return CountryInfo{IOCCode: countryCode, ISOCode: [2]string{countryCode}, Name: countryCode}
}

func GetCountryList() []CountryInfo {
	countries := make([]CountryInfo, 0, len(countriesData))
	for _, country := range countriesData {
		countries = append(countries, country)
	}

	return countries
}
