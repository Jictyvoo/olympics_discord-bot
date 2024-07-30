package entities

import (
	"crypto/sha256"
	"strconv"
	"time"
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
	StatusFinished  EventStatus = "finished"
)

type UnitType string

type OlympicCompetitors struct {
	ID          Identifier
	Name        string
	Code        string
	Age         uint8
	CountryCode string
	Country     CountryInfo
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

func (oe OlympicEvent) SHAIdentifier() string {
	hasher := sha256.New()
	hasher.Write([]byte(oe.DisciplineName))
	hasher.Write([]byte(strconv.Itoa(int(oe.Gender))))
	hasher.Write([]byte(oe.Phase))

	identifier := hasher.Sum([]byte(oe.EventName))
	return string(identifier)
}
