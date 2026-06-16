package fifafetch

import (
	"strconv"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const (
	sidesPerMatch = 2
	langKeyLen    = 2
	genderMale    = 1
	genderFemale  = 2
)

const (
	disciplineFootball = "Football"
	disciplineFutebol  = "Futebol"
)

// mapResults yields results only for a finished match with both scores present.
func mapResults(
	m apiMatch, fixtureID eventcore.CanonicalID, status eventcore.FixtureStatus,
) []eventcore.Result {
	if status != eventcore.FixtureFinished || m.HomeTeamScore == nil || m.AwayTeamScore == nil {
		return nil
	}
	out := make([]eventcore.Result, 0, sidesPerMatch)
	for _, side := range []struct {
		team  apiTeam
		score int
	}{
		{m.Home, *m.HomeTeamScore},
		{m.Away, *m.AwayTeamScore},
	} {
		if side.team.IdTeam == "" {
			continue
		}
		out = append(out, eventcore.Result{
			FixtureID:     fixtureID,
			ParticipantID: eventcore.NewID(eventcore.ProviderFIFA, side.team.IdTeam),
			Score:         strconv.Itoa(side.score),
			Outcome:       mapOutcome(m.Winner, side.team.IdTeam),
		})
	}
	return out
}

func mapStandings(resp apiStandingResponse) []eventcore.Standing {
	out := make([]eventcore.Standing, 0, len(resp.Results))
	for _, s := range resp.Results {
		out = append(out, eventcore.Standing{
			StageID:       eventcore.NewID(eventcore.ProviderFIFA, s.IdStage),
			ParticipantID: eventcore.NewID(eventcore.ProviderFIFA, s.Team.IdTeam),
			Rank:          s.Position,
			Points:        s.Points,
			Stats: map[string]any{
				"won":             s.Won,
				"lost":            s.Lost,
				"drawn":           s.Drawn,
				"played":          s.Played,
				"for":             s.For,
				"against":         s.Against,
				"goal_difference": s.GoalsDiference,
			},
		})
	}
	return out
}

func mapStatus(raw int) eventcore.FixtureStatus {
	switch raw {
	case statusFinished:
		return eventcore.FixtureFinished
	case statusLive:
		return eventcore.FixtureLive
	default:
		return eventcore.FixtureScheduled
	}
}

var footballDisciplineByLang = map[string]string{
	"pt": disciplineFutebol,
	"en": disciplineFootball,
	"es": "Fútbol",
	"fr": disciplineFootball,
	"de": "Fußball",
	"it": "Calcio",
}

// footballDiscipline is prefix tolerant ("pt-BR" -> "pt") and falls back to English.
func footballDiscipline(lang string) string {
	key := strings.ToLower(lang)
	if len(key) >= langKeyLen {
		key = key[:langKeyLen]
	}
	if label, ok := footballDisciplineByLang[key]; ok {
		return label
	}
	return disciplineFootball
}

// mapGender maps the upstream gender (1 male, 2 female) to the generic code.
func mapGender(g int) string {
	switch g {
	case genderMale:
		return "M"
	case genderFemale:
		return "F"
	}
	return ""
}

// mapOutcome treats an empty winner on a finished match as a draw.
func mapOutcome(winnerTeamID, teamID string) eventcore.Outcome {
	if winnerTeamID == "" {
		return eventcore.OutcomeDraw
	}
	if winnerTeamID == teamID {
		return eventcore.OutcomeWin
	}
	return eventcore.OutcomeLoss
}

// localized matches lang case-insensitively and prefix-tolerantly ("en" ~ "en-GB"),
// falling back to the first entry.
func localized(texts localizedText, lang string) string {
	if len(texts) == 0 {
		return ""
	}
	want := strings.ToLower(lang)
	for _, t := range texts {
		loc := strings.ToLower(t.Locale)
		if loc == want || strings.HasPrefix(loc, want) || strings.HasPrefix(want, loc) {
			return t.Description
		}
	}
	return texts[0].Description
}
