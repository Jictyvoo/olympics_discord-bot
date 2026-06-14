package eventcore

type ParticipantKind string

const (
	ParticipantTeam    ParticipantKind = "team"
	ParticipantAthlete ParticipantKind = "athlete"
)

func (k ParticipantKind) Valid() bool {
	return k == ParticipantTeam || k == ParticipantAthlete
}

type Participant struct {
	ID         CanonicalID
	Ext        ExternalID
	Kind       ParticipantKind
	Name       string
	Code       string
	CountryISO string
	Gender     string // provider-defined, not enum-restricted (e.g. "M", "F", "mixed")
}
