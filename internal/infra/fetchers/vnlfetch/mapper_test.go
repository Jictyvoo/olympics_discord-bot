package vnlfetch

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

// scoreSentinel is the int32-min value the feed uses for an unplayed match.
const scoreSentinel = -2147483648

const (
	compSlug     = "vnl-2026"
	roundSemana2 = "Semana 2"
	cityPasig    = "Pasig City"
)

func baseTeams() []apiTeam {
	return []apiTeam{
		{No: 1, Code: "AAA", Name: "Alpha"},
		{No: 2, Code: "BBB", Name: "Bravo"},
		{No: 3, Code: "CCC", Name: "Charlie"},
	}
}

//nolint:funlen // table-driven fixture
func TestMapSchedule_ResultsOnlyWhenFinished(t *testing.T) {
	winner := 2
	resp := apiScheduleResponse{
		AllTeams: baseTeams(),
		Matches: []apiMatch{
			{
				MatchNo:         10,
				MatchDateUtc:    mustTime("2026-06-20T18:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1661,
				Gender:          "Men",
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         2,
				TeamAScore:      scoreSentinel,
				TeamBScore:      scoreSentinel,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "Grupo 1"},
			},
			{
				MatchNo:         11,
				MatchDateUtc:    mustTime("2026-06-20T20:00:00Z"),
				MatchStatus:     statusFinished,
				TournamentNo:    1661,
				Gender:          "Men",
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         2,
				WinnerTeamNo:    &winner,
				TeamAScore:      1,
				TeamBScore:      3,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "Grupo 1"},
			},
		},
	}

	mapped := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T18:30:00Z"),
	)

	if len(mapped.fixtures) != 2 {
		t.Fatalf("expected 2 fixtures; got %d", len(mapped.fixtures))
	}
	// Only the finished match contributes results, one per side.
	if len(mapped.results) != 2 {
		t.Fatalf("expected 2 results; got %d", len(mapped.results))
	}
	winID := eventcore.NewID(eventcore.ProviderVNL, teamKey(2))
	for _, r := range mapped.results {
		want := eventcore.OutcomeLoss
		if r.ParticipantID == winID {
			want = eventcore.OutcomeWin
		}
		if r.Outcome != want {
			t.Errorf("participant outcome = %q, want %q", r.Outcome, want)
		}
	}
	// Two teams deduped, one pool, one stage.
	if len(mapped.participants) != 2 {
		t.Fatalf("expected 2 participants; got %d", len(mapped.participants))
	}
	if len(mapped.groups) != 1 || len(mapped.stages) != 1 {
		t.Fatalf(
			"expected 1 group and 1 stage; got %d / %d",
			len(mapped.groups),
			len(mapped.stages),
		)
	}
}

//nolint:funlen // table-driven fixture
func TestMapSchedule_CompositeStageKeysPreventCollision(t *testing.T) {
	// Same roundNo (297) in two tournaments must yield two distinct stages, while
	// the same tournament+round dedups to one.
	resp := apiScheduleResponse{
		AllTeams: baseTeams(),
		Matches: []apiMatch{
			{
				MatchNo:         1,
				MatchDateUtc:    mustTime("2026-06-20T10:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1661,
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         2,
				RoundNo:         297,
				RoundName:       roundSemana2,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
			},
			{
				MatchNo:         2,
				MatchDateUtc:    mustTime("2026-06-20T12:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1662,
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         3,
				RoundNo:         297,
				RoundName:       roundSemana2,
				RoundCode:       "2",
				Pool:            apiPool{No: 200, Name: "G2"},
			},
			{
				MatchNo:         3,
				MatchDateUtc:    mustTime("2026-06-20T14:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1661,
				CompetitionSlug: compSlug,
				TeamANo:         2,
				TeamBNo:         3,
				RoundNo:         297,
				RoundName:       roundSemana2,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
			},
		},
	}

	mapped := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T00:00:00Z"),
	)

	if len(mapped.competitions) != 1 {
		t.Fatalf("expected 1 shared competition; got %d", len(mapped.competitions))
	}
	if len(mapped.seasons) != 2 {
		t.Fatalf("expected 2 seasons; got %d", len(mapped.seasons))
	}
	if len(mapped.stages) != 2 {
		t.Fatalf("expected 2 stages (one per tournament); got %d", len(mapped.stages))
	}
	if len(mapped.groups) != 2 {
		t.Fatalf("expected 2 groups; got %d", len(mapped.groups))
	}
}

//nolint:funlen // table-driven fixture
func TestMapSchedule_EmitsVenueFromCity(t *testing.T) {
	resp := apiScheduleResponse{
		AllTeams: baseTeams(),
		Matches: []apiMatch{
			{
				MatchNo:         1,
				MatchDateUtc:    mustTime("2026-06-20T10:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1662,
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         2,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
				City:            cityPasig,
				CountryCode:     "PH",
			},
			{
				// Same host city, deduped to one venue.
				MatchNo:         2,
				MatchDateUtc:    mustTime("2026-06-20T12:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1662,
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         3,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
				City:            cityPasig,
				CountryCode:     "PH",
			},
			{
				// No city -> no venue, fixture carries a nil venue id.
				MatchNo:         3,
				MatchDateUtc:    mustTime("2026-06-20T14:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1662,
				CompetitionSlug: compSlug,
				TeamANo:         2,
				TeamBNo:         3,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
			},
		},
	}

	mapped := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T00:00:00Z"),
	)

	if len(mapped.venues) != 1 {
		t.Fatalf("expected 1 deduped venue; got %d", len(mapped.venues))
	}
	v := mapped.venues[0]
	if v.City != cityPasig || v.CountryISO != "PH" {
		t.Errorf("venue = %+v; want city=Pasig City country=PH", v)
	}
	var withCity, withoutCity int
	for _, f := range mapped.fixtures {
		if f.VenueID != nil {
			withCity++
		} else {
			withoutCity++
		}
	}
	if withCity != 2 || withoutCity != 1 {
		t.Errorf("venue ids: with=%d without=%d; want 2/1", withCity, withoutCity)
	}
}

func TestMapSchedule_CompletesLiveAfterGrace(t *testing.T) {
	winner := 1
	match := apiMatch{
		MatchNo: 50, MatchDateUtc: mustTime("2026-06-20T18:00:00Z"), MatchStatus: statusLive,
		TournamentNo: 1661, CompetitionSlug: compSlug, TeamANo: 1, TeamBNo: 2,
		WinnerTeamNo: &winner, TeamAScore: 3, TeamBScore: 1,
		RoundNo: 297, RoundCode: "2", Pool: apiPool{No: 100, Name: "G1"},
	}
	resp := apiScheduleResponse{AllTeams: baseTeams(), Matches: []apiMatch{match}}

	// Ends at 20:00; grace expires at 21:00.
	withinGrace := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T20:30:00Z"),
	)
	if got := withinGrace.fixtures[0].Status; got != eventcore.FixtureLive {
		t.Errorf("within grace: status = %q; want live", got)
	}
	if len(withinGrace.results) != 0 {
		t.Errorf("within grace: expected no results; got %d", len(withinGrace.results))
	}

	afterGrace := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T21:30:00Z"),
	)
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

func TestMapSchedule_SkipsPlaceholderFixtures(t *testing.T) {
	resp := apiScheduleResponse{
		AllTeams: baseTeams(),
		Matches: []apiMatch{
			{
				// Finals slot: teams not drawn yet.
				MatchNo:         99,
				MatchDateUtc:    mustTime("2026-07-30T18:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1661,
				CompetitionSlug: compSlug,
				IsMatchTBD:      true,
				RoundNo:         300,
				RoundCode:       "3",
			},
			{
				MatchNo:         7,
				MatchDateUtc:    mustTime("2026-06-20T18:00:00Z"),
				MatchStatus:     statusScheduled,
				TournamentNo:    1661,
				CompetitionSlug: compSlug,
				TeamANo:         1,
				TeamBNo:         2,
				RoundNo:         297,
				RoundCode:       "2",
				Pool:            apiPool{No: 100, Name: "G1"},
			},
		},
	}

	mapped := mapSchedule(
		resp,
		"en",
		mustTime("2026-06-20T00:00:00Z"),
		mustTime("2026-06-20T18:30:00Z"),
	)

	if len(mapped.fixtures) != 1 {
		t.Fatalf("expected only the well-defined fixture; got %d", len(mapped.fixtures))
	}
	if mapped.fixtures[0].Ext.Key != matchKey(7) {
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
