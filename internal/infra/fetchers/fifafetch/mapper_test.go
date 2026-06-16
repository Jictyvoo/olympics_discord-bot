package fifafetch

import (
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func mustTime(s string) time.Time {
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return ts
}

func TestMapMatches_DrawNoResultsUntilFinished(t *testing.T) {
	resp := apiMatchesResponse{Results: []apiMatch{
		{
			IdMatch:       "m-scheduled",
			IdCompetition: "17",
			IdSeason:      "285023",
			IdStage:       "s1",
			Date:          mustTime("2026-06-20T18:00:00Z"),
			Home:          apiTeam{IdTeam: "A", Abbreviation: "A", Gender: 1},
			Away:          apiTeam{IdTeam: "B", Abbreviation: "B", Gender: 1},
			MatchStatus:   statusScheduled,
		},
		{
			IdMatch:       "m-draw",
			IdCompetition: "17",
			IdSeason:      "285023",
			IdStage:       "s1",
			Date:          mustTime("2026-06-20T20:00:00Z"),
			Home:          apiTeam{IdTeam: "A", Abbreviation: "A", Gender: 1},
			Away:          apiTeam{IdTeam: "C", Abbreviation: "C", Gender: 1},
			HomeTeamScore: new(1),
			AwayTeamScore: new(1),
			Winner:        "",
			MatchStatus:   statusFinished,
		},
	}}

	mapped := mapMatches(resp, "en", seasonMeta{})

	if len(mapped.fixtures) != 2 {
		t.Fatalf("expected 2 fixtures; got %d", len(mapped.fixtures))
	}
	// Scheduled match yields no results; finished draw yields 2 draw results.
	if len(mapped.results) != 2 {
		t.Fatalf("expected 2 results (draw only); got %d", len(mapped.results))
	}
	for _, r := range mapped.results {
		if r.Outcome != eventcore.OutcomeDraw {
			t.Errorf("expected draw outcome, got %q", r.Outcome)
		}
	}
	// Three distinct teams (A, B, C) deduped across both matches.
	if len(mapped.participants) != 3 {
		t.Fatalf("expected 3 participants; got %d", len(mapped.participants))
	}
	if len(mapped.stages) != 1 || len(mapped.stageKeys) != 1 {
		t.Fatalf("expected 1 deduped stage; got %d / %v", len(mapped.stages), mapped.stageKeys)
	}
}

func TestMapMatches_SkipsPlaceholderFixtures(t *testing.T) {
	resp := apiMatchesResponse{Results: []apiMatch{
		{
			// Knockout slot: no teams yet.
			IdMatch:       "placeholder",
			IdCompetition: "17",
			IdSeason:      "285023",
			IdStage:       "knockout",
			Date:          mustTime("2026-07-04T18:00:00Z"),
			MatchStatus:   statusScheduled,
		},
		{
			IdMatch:       "real",
			IdCompetition: "17",
			IdSeason:      "285023",
			IdStage:       "groups",
			Date:          mustTime("2026-06-20T18:00:00Z"),
			Home:          apiTeam{IdTeam: "A", Abbreviation: "A", Gender: 1},
			Away:          apiTeam{IdTeam: "B", Abbreviation: "B", Gender: 1},
			MatchStatus:   statusScheduled,
		},
	}}

	mapped := mapMatches(resp, "en", seasonMeta{})

	if len(mapped.fixtures) != 1 {
		t.Fatalf("expected only the well-defined fixture; got %d", len(mapped.fixtures))
	}
	if mapped.fixtures[0].Ext.Key != "real" {
		t.Errorf("kept the wrong fixture: %q", mapped.fixtures[0].Ext.Key)
	}
	if len(mapped.stages) != 1 || len(mapped.participants) != 2 {
		t.Errorf(
			"placeholder leaked: stages=%d participants=%d",
			len(mapped.stages),
			len(mapped.participants),
		)
	}
}
