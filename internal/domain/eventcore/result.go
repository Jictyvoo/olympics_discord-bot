package eventcore

type Outcome string

const (
	OutcomeNone        Outcome = ""
	OutcomeWin         Outcome = "win"
	OutcomeLoss        Outcome = "loss"
	OutcomeDraw        Outcome = "draw"
	OutcomeMedalGold   Outcome = "medal_gold"
	OutcomeMedalSilver Outcome = "medal_silver"
	OutcomeMedalBronze Outcome = "medal_bronze"
)

func (o Outcome) Valid() bool {
	switch o {
	case OutcomeNone, OutcomeWin, OutcomeLoss, OutcomeDraw,
		OutcomeMedalGold, OutcomeMedalSilver, OutcomeMedalBronze:
		return true
	}
	return false
}

type Result struct {
	FixtureID     CanonicalID
	ParticipantID CanonicalID
	Position      *int
	Score         string
	RawMark       string
	Outcome       Outcome
}

type Standing struct {
	StageID       CanonicalID
	ParticipantID CanonicalID
	Rank          int
	Points        int
	Stats         map[string]any
}
