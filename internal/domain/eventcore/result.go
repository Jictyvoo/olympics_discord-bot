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

func (m medalType) value() int {
	switch m {
	case MedalBronze:
		return 3
	case MedalSilver:
		return 5
	case MedalGold:
		return 7
	case MedalLoser:
		return 1
	case MedalWinner:
		return 2
	}

	return 0
}

func (m medalType) CompareTo(other medalType) int {
	result := m.value() - other.value()
	if result < 0 {
		return -1
	} else if result > 1 {
		return 1
	}

	return 0
}

type Results struct {
	Position  string
	Mark      string
	MedalType medalType
	Irm       string
}
