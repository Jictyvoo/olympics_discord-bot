package eventcore

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type FixtureStatus string

const (
	FixtureScheduled FixtureStatus = "scheduled"
	FixtureLive      FixtureStatus = "live"
	FixtureFinished  FixtureStatus = "finished"
	FixtureCancelled FixtureStatus = "cancelled"
	FixturePostponed FixtureStatus = "postponed"
)

func (s FixtureStatus) Valid() bool {
	switch s {
	case FixtureScheduled, FixtureLive, FixtureFinished, FixtureCancelled, FixturePostponed:
		return true
	}
	return false
}

type FixtureParticipant struct {
	ParticipantID CanonicalID
	Role          string // provider-defined, e.g. "home", "away", "athlete"
}

type Fixture struct {
	ID           CanonicalID
	Ext          ExternalID
	StageID      CanonicalID
	GroupID      *CanonicalID
	VenueID      *CanonicalID
	Name         string
	StartsAt     time.Time
	EndsAt       time.Time
	Status       FixtureStatus
	Checksum     string // SHA-256 of normalised payload; gates re-notification
	Participants []FixtureParticipant
}

// ComputeChecksum returns the fixture-only checksum (no results). Equivalent to
// ComputeChecksumWith(nil); kept for call sites that have no results to hash.
func (f Fixture) ComputeChecksum() string {
	return f.ComputeChecksumWith(nil)
}

// ComputeChecksumWith returns a stable hex-encoded SHA-256 over the fixture's
// fields and its results: the data that, when changed, should trigger a
// re-notification or re-sync. Participants and results are sorted first so
// upstream ordering can never churn the checksum. The receiver is not modified.
func (f Fixture) ComputeChecksumWith(results []Result) string {
	type stableResult struct {
		ParticipantID string
		Position      *int
		Score         string
		RawMark       string
		Outcome       Outcome
	}
	type stable struct {
		ExtKey   string
		Name     string
		StartsAt time.Time
		EndsAt   time.Time
		Status   FixtureStatus
		Parts    []FixtureParticipant
		Results  []stableResult
	}

	parts := make([]FixtureParticipant, len(f.Participants))
	copy(parts, f.Participants)
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].ParticipantID.String() < parts[j].ParticipantID.String()
	})

	stableResults := make([]stableResult, 0, len(results))
	for _, r := range results {
		stableResults = append(stableResults, stableResult{
			ParticipantID: r.ParticipantID.String(),
			Position:      r.Position,
			Score:         r.Score,
			RawMark:       r.RawMark,
			Outcome:       r.Outcome,
		})
	}
	sort.Slice(stableResults, func(i, j int) bool {
		return stableResults[i].ParticipantID < stableResults[j].ParticipantID
	})

	payload, _ := json.Marshal(stable{
		ExtKey:   f.Ext.Key,
		Name:     f.Name,
		StartsAt: f.StartsAt.UTC(),
		EndsAt:   f.EndsAt.UTC(),
		Status:   f.Status,
		Parts:    parts,
		Results:  stableResults,
	})
	sum := sha256.Sum256(payload)
	return fmt.Sprintf("%x", sum)
}
