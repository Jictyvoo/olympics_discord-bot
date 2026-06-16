package fifafetch

import (
	"testing"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func TestMapStatus(t *testing.T) {
	cases := []struct {
		raw  int
		want eventcore.FixtureStatus
	}{
		{statusFinished, eventcore.FixtureFinished},
		{statusScheduled, eventcore.FixtureScheduled},
		{statusLive, eventcore.FixtureLive},
		{99, eventcore.FixtureScheduled},
	}
	for _, tc := range cases {
		if got := mapStatus(tc.raw); got != tc.want {
			t.Errorf("mapStatus(%d) = %q, want %q", tc.raw, got, tc.want)
		}
	}
}

func TestMapOutcome(t *testing.T) {
	cases := []struct {
		name         string
		winner, team string
		want         eventcore.Outcome
	}{
		{"win", "43911", "43911", eventcore.OutcomeWin},
		{"loss", "43911", "43883", eventcore.OutcomeLoss},
		{"draw", "", "43911", eventcore.OutcomeDraw},
	}
	for _, tc := range cases {
		if got := mapOutcome(tc.winner, tc.team); got != tc.want {
			t.Errorf(
				"%s: mapOutcome(%q,%q) = %q, want %q",
				tc.name,
				tc.winner,
				tc.team,
				got,
				tc.want,
			)
		}
	}
}

func TestFootballDiscipline(t *testing.T) {
	cases := []struct {
		lang, want string
	}{
		{"pt", disciplineFutebol},
		{"pt-BR", disciplineFutebol},
		{"EN", disciplineFootball},
		{"es", "Fútbol"},
		{"ja", disciplineFootball},
		{"", disciplineFootball},
	}
	for _, tc := range cases {
		if got := footballDiscipline(tc.lang); got != tc.want {
			t.Errorf("footballDiscipline(%q) = %q, want %q", tc.lang, got, tc.want)
		}
	}
}

func TestLocalized(t *testing.T) {
	texts := localizedText{
		{Locale: "fr-FR", Description: "Mexique"},
		{Locale: "en-GB", Description: "Mexico"},
	}
	if got := localized(texts, "en"); got != "Mexico" {
		t.Errorf("localized(en) = %q, want Mexico", got)
	}
	if got := localized(nil, "en"); got != "" {
		t.Errorf("localized(nil) = %q, want empty", got)
	}
	if got := localized(localizedText{{Locale: "pt-BR", Description: "x"}}, "en"); got != "x" {
		t.Errorf("localized fallback = %q, want x", got)
	}
}

func TestMapStandings(t *testing.T) {
	resp := apiStandingResponse{Results: []apiStanding{
		{
			IdStage: "s1", Position: 1, Points: 3, Won: 1, GoalsDiference: 2,
			Team: apiStandingTeam{IdTeam: "A"},
		},
	}}
	got := mapStandings(resp)
	if len(got) != 1 {
		t.Fatalf("expected 1 standing; got %d", len(got))
	}
	if got[0].Rank != 1 || got[0].Points != 3 {
		t.Errorf("unexpected standing %+v", got[0])
	}
	if got[0].Stats["goal_difference"] != 2 {
		t.Errorf("expected goal_difference 2, got %v", got[0].Stats["goal_difference"])
	}
}
