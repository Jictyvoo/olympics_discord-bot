package entities

import "strings"

type medalType string

const (
	MedalNoMedal medalType = ""
	MedalBronze  medalType = "bronze"
	MedalSilver  medalType = "silver"
	MedalGold    medalType = "gold"
	MedalWinner  medalType = "winner"
	MedalLoser   medalType = "loser"
)

//goland:noinspection GoExportedFuncWithUnexportedType
func Medal(value string) medalType {
	value = strings.ToLower(value)
	switch value {
	case "w":
		return MedalWinner
	case "l":
		return MedalLoser
	}

	switch {
	case strings.Contains(value, "bronze"):
		return MedalBronze
	case strings.Contains(value, "silver"):
		return MedalSilver
	case strings.Contains(value, "gold"):
		return MedalGold
	}

	return MedalNoMedal
}

type Results struct {
	Position  string
	Mark      string
	MedalType medalType
	Irm       string
}
