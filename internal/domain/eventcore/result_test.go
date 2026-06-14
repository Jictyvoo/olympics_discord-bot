package eventcore

import "testing"

func TestOutcome_Valid(t *testing.T) {
	testCases := []struct {
		name    string
		outcome Outcome
		want    bool
	}{
		{"none", OutcomeNone, true},
		{"win", OutcomeWin, true},
		{"loss", OutcomeLoss, true},
		{"draw", OutcomeDraw, true},
		{"medal_gold", OutcomeMedalGold, true},
		{"medal_silver", OutcomeMedalSilver, true},
		{"medal_bronze", OutcomeMedalBronze, true},
		{unknownProvider, "unknown_outcome", false},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			if got := tCase.outcome.Valid(); got != tCase.want {
				t.Fatalf("Valid() = %v, want %v", got, tCase.want)
			}
		})
	}
}
