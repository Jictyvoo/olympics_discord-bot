package entities

type OlympicCompetitors struct {
	ID          Identifier
	Name        string
	Code        string
	Age         uint8
	CountryCode string
	Country     CountryInfo
}
