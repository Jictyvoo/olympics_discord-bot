package entities

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jictyvoo/olympics_data_fetcher/internal/utils"
)

type (
	Identifier uint64
	HexID      string
)

type Gender uint8

const (
	GenderOther Gender = iota
	GenderMasc
	GenderFem
)

func (g Gender) String() string {
	switch g {
	case GenderMasc:
		return "Male"
	case GenderFem:
		return "Female"
	default:
		return "Other"
	}
}

type EventStatus string

const (
	StatusScheduled EventStatus = "scheduled"
	StatusOngoing   EventStatus = "ongoing"
	StatusFinished  EventStatus = "finished"
)

type UnitType string

type OlympicEvent struct {
	ID                  Identifier
	EventName           string
	Discipline          Discipline
	Phase               string
	Gender              Gender
	SessionCode         string
	UnitType            UnitType
	StartAt             time.Time
	EndAt               time.Time
	Status              EventStatus
	HasMedal            bool
	Competitors         []OlympicCompetitors
	ResultPerCompetitor map[string]Results
}

//goland:noinspection GoMixedReceiverTypes
func (oe OlympicEvent) SHAIdentifier() string {
	competitorsResults := make([]utils.KeyValueEntry[Results], len(oe.ResultPerCompetitor))
	var index int
	for code, compResult := range oe.ResultPerCompetitor {
		competitorsResults[index] = utils.KeyValueEntry[Results]{
			Key:   code,
			Value: compResult,
		}
		index++
	}

	oe.ResultPerCompetitor = nil
	var comparator struct {
		OlympicEvent
		ResultsPerCompetitor []utils.KeyValueEntry[Results]
	}
	comparator.OlympicEvent = oe
	comparator.ResultsPerCompetitor = competitorsResults
	slices.SortFunc(comparator.ResultsPerCompetitor, utils.KeyValueEntry[Results].Compare)
	if hash, err := utils.Hash(comparator); err == nil && len(hash) > 0 {
		return hash
	}

	return ""
}

//goland:noinspection GoMixedReceiverTypes
func (oe *OlympicEvent) Normalize() {
	oe.StartAt = oe.StartAt.In(time.UTC)
	oe.EndAt = oe.EndAt.In(time.UTC)
	slices.SortFunc(
		oe.Competitors, func(a, b OlympicCompetitors) int {
			var results struct{ a, b Results }
			results.a = oe.ResultPerCompetitor[a.Code]
			results.b = oe.ResultPerCompetitor[b.Code]
			if results.a.MedalType != MedalNoMedal || results.b.MedalType != MedalNoMedal {
				return results.b.MedalType.CompareTo(results.a.MedalType)
			}
			if results.a.Mark != "" || results.b.Mark != "" {
				return compareMark(results.b.Mark, results.a.Mark)
			}

			return strings.Compare(a.Code, b.Code)
		},
	)
}

func compareMark(a, b string) int {
	// Try to compare as float
	var (
		fA, fB float64
		err    [2]error
	)
	fA, err[0] = strconv.ParseFloat(a, 64)
	fB, err[1] = strconv.ParseFloat(b, 64)
	if err[0] == nil && err[1] == nil {
		return int(fA - fB)
	}
	return strings.Compare(a, b)
}
