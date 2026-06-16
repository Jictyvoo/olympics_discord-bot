package eventcore

// FixtureContext is a read aggregate locating a fixture (competition, stage and
// optional group names), resolved in a single query for rendering.
type FixtureContext struct {
	Competition Competition
	StageName   string
	StageOrd    int
	GroupName   string
}

// FixtureCompetitor is a read aggregate: a fixture participant with its role,
// resolved ISO2 code, and result (zero-value when there is no result yet).
type FixtureCompetitor struct {
	Participant Participant
	Role        string
	CountryISO2 string
	Result      Result
}
