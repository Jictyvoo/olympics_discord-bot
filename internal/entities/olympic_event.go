package entities

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"hash/fnv"
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
	ID          Identifier
	EventName   string
	Discipline  Discipline
	Phase       string
	Gender      Gender
	SessionCode string
	UnitType    UnitType
	StartAt     time.Time
	EndAt       time.Time
	Status      EventStatus
	Competitors []OlympicCompetitors
}

func (oe OlympicEvent) Hash() (string, error) {
	h := fnv.New64a() // Create a new FNV-1a 64-bit hash instance

	// Serialize POSMessageConfig into gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(oe); err != nil {
		return "", err
	}

	// Write the serialized bytes to the hash
	if _, err := h.Write(buf.Bytes()); err != nil {
		return "", err
	}

	// Return the resulting hash value as a hexadecimal string
	return fmt.Sprintf("%x", h.Sum64()), nil
}

func (oe OlympicEvent) SHAIdentifier() string {
	if hash, err := oe.Hash(); err == nil && len(hash) > 0 {
		return hash
	}

	var buffer bytes.Buffer
	buffer.Write([]byte(oe.Discipline.Name))
	buffer.Write([]byte(strconv.Itoa(int(oe.Gender))))
	buffer.Write([]byte(oe.Phase))
	buffer.Write([]byte(oe.SessionCode))

	identifier := sha256.Sum256([]byte(oe.EventName))
	return string(identifier[:])
}
