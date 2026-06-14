package eventcore

import "errors"

// ErrNotImplemented is returned by provider methods that have not been implemented yet.
var ErrNotImplemented = errors.New("provider: not implemented")

// SyncDelta carries all records produced by a single provider sync operation.
type SyncDelta struct {
	Competitions []Competition
	Seasons      []Season
	Stages       []Stage
	Groups       []Group
	Venues       []Venue
	Participants []Participant
	Fixtures     []Fixture
	Results      []Result
	Standings    []Standing
	Cursor       string
}
