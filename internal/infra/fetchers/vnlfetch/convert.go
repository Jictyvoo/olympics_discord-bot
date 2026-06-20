package vnlfetch

import (
	"strconv"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func mapResults(
	m apiMatch, fixtureID eventcore.CanonicalID, status eventcore.FixtureStatus,
) []eventcore.Result {
	if status != eventcore.FixtureFinished || m.TeamAScore < 0 || m.TeamBScore < 0 {
		return nil
	}
	out := make([]eventcore.Result, 0, sidesPerMatch)
	for _, side := range []struct {
		no    int
		score int
	}{
		{m.TeamANo, m.TeamAScore},
		{m.TeamBNo, m.TeamBScore},
	} {
		if side.no == 0 {
			continue
		}
		out = append(out, eventcore.Result{
			FixtureID:     fixtureID,
			ParticipantID: eventcore.NewID(eventcore.ProviderVNL, teamKey(side.no)),
			Score:         strconv.Itoa(side.score),
			Outcome:       mapOutcome(m.WinnerTeamNo, side.no),
		})
	}
	return out
}

func mapOutcome(winnerTeamNo *int, teamNo int) eventcore.Outcome {
	if winnerTeamNo == nil {
		return eventcore.OutcomeNone
	}
	if *winnerTeamNo == teamNo {
		return eventcore.OutcomeWin
	}
	return eventcore.OutcomeLoss
}

func mapGender(g string) string {
	switch strings.ToLower(g) {
	case "men", "male", "m":
		return "M"
	case "women", "female", "w":
		return "F"
	}
	return ""
}

const disciplineVolleyball = "Volleyball"

var volleyballDisciplineByLang = map[string]string{
	"pt": "Vôlei",
	"en": disciplineVolleyball,
	"es": "Voleibol",
	"fr": "Volley-ball",
	"de": disciplineVolleyball,
	"it": "Pallavolo",
}

func volleyballDiscipline(lang string) string {
	key := strings.ToLower(lang)
	if len(key) >= langKeyLen {
		key = key[:langKeyLen]
	}
	if label, ok := volleyballDisciplineByLang[key]; ok {
		return label
	}
	return disciplineVolleyball
}
