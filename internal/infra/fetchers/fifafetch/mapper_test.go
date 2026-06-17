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

	mapped := mapMatches(resp, "en", seasonMeta{}, mustTime("2026-06-20T18:30:00Z"))

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
	// The feed has no kickoff end; fixtures end 105 min after the start.
	for _, f := range mapped.fixtures {
		if want := f.StartsAt.Add(footballMatchDuration); !f.EndsAt.Equal(want) {
			t.Errorf("fixture %q EndsAt = %s; want %s", f.Ext.Key, f.EndsAt, want)
		}
	}
}

func TestMapMatches_CompletesLiveAfterGrace(t *testing.T) {
	match := apiMatch{
		IdMatch:       "m-live",
		IdCompetition: "17",
		IdSeason:      "285023",
		IdStage:       "s1",
		Date:          mustTime("2026-06-18T22:00:00Z"),
		Home:          apiTeam{IdTeam: "CA", Abbreviation: "CA", Gender: 1},
		Away:          apiTeam{IdTeam: "QA", Abbreviation: "QA", Gender: 1},
		HomeTeamScore: new(2),
		AwayTeamScore: new(1),
		Winner:        "CA",
		MatchStatus:   statusLive,
	}
	resp := apiMatchesResponse{Results: []apiMatch{match}}

	// Ends at 23:45; grace expires at 00:45 next day.
	withinGrace := mapMatches(resp, "en", seasonMeta{}, mustTime("2026-06-19T00:30:00Z"))
	if len(withinGrace.fixtures) != 1 {
		t.Fatalf("expected 1 fixture; got %d", len(withinGrace.fixtures))
	}
	if got := withinGrace.fixtures[0].Status; got != eventcore.FixtureLive {
		t.Errorf("within grace: status = %q; want live", got)
	}
	if len(withinGrace.results) != 0 {
		t.Errorf("within grace: expected no results; got %d", len(withinGrace.results))
	}

	afterGrace := mapMatches(resp, "en", seasonMeta{}, mustTime("2026-06-19T01:00:00Z"))
	if got := afterGrace.fixtures[0].Status; got != eventcore.FixtureFinished {
		t.Errorf("after grace: status = %q; want finished", got)
	}
	if len(afterGrace.results) != 2 {
		t.Errorf("after grace: expected 2 results; got %d", len(afterGrace.results))
	}
	if withinGrace.fixtures[0].Checksum == afterGrace.fixtures[0].Checksum {
		t.Error("checksum did not change when the fixture completed")
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

	mapped := mapMatches(resp, "en", seasonMeta{}, mustTime("2026-06-20T18:30:00Z"))

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
