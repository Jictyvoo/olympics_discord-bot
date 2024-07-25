package entities

import "time"

type Gender uint8

const (
	GenderOther Gender = iota
	GenderMasc
	GenderFem
)

type EventStatus string

const (
	StatusScheduled EventStatus = "scheduled"
	StatusFinished  EventStatus = "finished"
)

type UnitType string

type Identifier uint64

type OlympicCompetitors struct {
	ID          Identifier
	Code        string
	CountryCode string
	Name        string
	CountryInfo CountryInfo
}

type OlympicEvent struct {
	ID             Identifier
	EventName      string
	DisciplineName string
	Phase          string
	Gender         Gender
	UnitType       UnitType
	StartAt        time.Time
	EndAt          time.Time
	Status         EventStatus
	Competitors    []OlympicCompetitors
}
