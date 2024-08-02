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
	case "w", string(MedalWinner):
		return MedalWinner
	case "l", string(MedalLoser):
		return MedalLoser
	case "b", string(MedalBronze):
		return MedalBronze
	case "s", string(MedalSilver):
		return MedalSilver
	case "g", string(MedalGold):
		return MedalGold
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

func (m medalType) String() string {
	switch m {
	case MedalBronze:
		return ":third_place:"
	case MedalSilver:
		return ":second_place:"
	case MedalGold:
		return ":first_place:"
	}

	return string(m)
}

type Results struct {
	Position  string
	Mark      string
	MedalType medalType
	Irm       string
}
