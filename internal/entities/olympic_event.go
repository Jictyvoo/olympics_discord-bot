package entities

import (
	"slices"
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
