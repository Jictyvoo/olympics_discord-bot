package fifafetch

import (
	"strconv"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

// mapResults yields one result per side, only for a finished match with both
// scores present.
func mapResults(
	m apiMatch, fixtureID eventcore.CanonicalID, status eventcore.FixtureStatus,
) []eventcore.Result {
	if status != eventcore.FixtureFinished || m.HomeTeamScore == nil || m.AwayTeamScore == nil {
		return nil
	}
	out := make([]eventcore.Result, 0, 2)
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

// mapGender maps the upstream gender (1 male, 2 female) to the generic code.
func mapGender(g int) string {
	switch g {
	case 1:
		return "M"
	case 2:
		return "F"
	}
	return ""
}

// mapOutcome derives a side's outcome from the winner team id; empty winner on a
// finished match is a draw.
func mapOutcome(winnerTeamID, teamID string) eventcore.Outcome {
	if winnerTeamID == "" {
		return eventcore.OutcomeDraw
	}
	if winnerTeamID == teamID {
		return eventcore.OutcomeWin
	}
	return eventcore.OutcomeLoss
}

// localized returns the description matching lang (case-insensitive, prefix
// tolerant so "en" matches "en-GB"), falling back to the first entry.
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
