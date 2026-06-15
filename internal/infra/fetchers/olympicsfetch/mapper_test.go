package olympicsfetch

import (
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func sampleUnit() apiUnit {
	u := apiUnit{
		Id:             "u-100m-final",
		DisciplineName: "Athletics",
		DisciplineCode: "ATH",
		EventUnitName:  "100m Final",
		PhaseId:        "phase-final",
		PhaseName:      "Final",
		EventOrder:     3,
		GenderCode:     "W",
		Status:         "Finished",
		StartDate:      time.Date(2024, 8, 4, 20, 55, 0, 0, time.UTC),
		EndDate:        time.Date(2024, 8, 4, 21, 5, 0, 0, time.UTC),
		Competitors: []apiCompetitor{
			{Code: "USA-001", Noc: "USA", Name: "A. Sprinter"},
			{Code: "JAM-002", Noc: "JAM", Name: "B. Quickstep"},
		},
	}
	u.Competitors[0].Results.Position = "1"
	u.Competitors[0].Results.Mark = "9.79"
	u.Competitors[0].Results.MedalType = medalCodeGold
	u.Competitors[1].Results.Position = "2"
	u.Competitors[1].Results.Mark = "9.81"
	u.Competitors[1].Results.MedalType = medalCodeSilver
	return u
}

func TestMapSchedule_RelationalChain(t *testing.T) {
	out := mapSchedule(apiScheduleResponse{Units: []apiUnit{sampleUnit()}})

	if len(out.competitions) != 1 {
		t.Fatalf("competitions = %d, want 1", len(out.competitions))
	}
	if out.competitions[0].Code != "ATH" || out.competitions[0].Discipline != "Athletics" {
		t.Errorf("competition = %+v, want Code=ATH Discipline=Athletics", out.competitions[0])
	}
	if len(out.seasons) != 1 {
		t.Fatalf("seasons = %d, want 1", len(out.seasons))
	}
	if out.seasons[0].CompetitionID != out.competitions[0].ID {
		t.Errorf("season.CompetitionID does not link to emitted competition")
	}
	if len(out.stages) != 1 {
		t.Fatalf("stages = %d, want 1", len(out.stages))
	}
	if out.stages[0].Name != "Final" || out.stages[0].Ord != 3 {
		t.Errorf("stage = %+v, want Name=Final Ord=3", out.stages[0])
	}
	if out.stages[0].SeasonID != out.seasons[0].ID {
		t.Errorf("stage.SeasonID does not link to emitted season")
	}
	if len(out.fixtures) != 1 {
		t.Fatalf("fixtures = %d, want 1", len(out.fixtures))
	}
	if out.fixtures[0].StageID.IsZero() {
		t.Errorf("fixture.StageID is zero; must always be set")
	}
	if out.fixtures[0].StageID != out.stages[0].ID {
		t.Errorf("fixture.StageID does not match emitted stage")
	}
}

func TestMapSchedule_Dedupe(t *testing.T) {
	u1 := sampleUnit()
	u2 := sampleUnit()
	u2.Id = "u-100m-semi" // different fixture, same discipline + phase
	out := mapSchedule(apiScheduleResponse{Units: []apiUnit{u1, u2}})

	if len(out.competitions) != 1 {
		t.Errorf("competitions = %d, want 1 (deduped)", len(out.competitions))
	}
	if len(out.seasons) != 1 {
		t.Errorf("seasons = %d, want 1 (deduped)", len(out.seasons))
	}
	if len(out.stages) != 1 {
		t.Errorf("stages = %d, want 1 (deduped)", len(out.stages))
	}
	if len(out.participants) != 2 {
		t.Errorf("participants = %d, want 2 (deduped)", len(out.participants))
	}
	if len(out.fixtures) != 2 {
		t.Errorf("fixtures = %d, want 2", len(out.fixtures))
	}
}

func TestMapSchedule_Group(t *testing.T) {
	u := sampleUnit()
	u.GroupId = "grp-A"
	out := mapSchedule(apiScheduleResponse{
		Units:  []apiUnit{u},
		Groups: []apiGroup{{Id: "grp-A", Title: "Group A"}},
	})
	if len(out.groups) != 1 {
		t.Fatalf("groups = %d, want 1", len(out.groups))
	}
	if out.groups[0].StageID != out.stages[0].ID {
		t.Errorf("group.StageID does not link to emitted stage")
	}
	if out.fixtures[0].GroupID == nil || *out.fixtures[0].GroupID != out.groups[0].ID {
		t.Errorf("fixture.GroupID not set to emitted group")
	}
}

func TestMapGender(t *testing.T) {
	cases := []struct{ in, want string }{
		{"W", "F"},
		{"M", "M"},
		{"X", ""},
		{"", ""},
	}
	for _, tc := range cases {
		if got := mapGender(tc.in); got != tc.want {
			t.Errorf("mapGender(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMapSchedule_ParticipantGender(t *testing.T) {
	out := mapSchedule(apiScheduleResponse{Units: []apiUnit{sampleUnit()}})
	for _, p := range out.participants {
		if p.Gender != "F" {
			t.Errorf("participant %s gender = %q, want F", p.Code, p.Gender)
		}
	}
}

func TestMapResults_EmptyFilter(t *testing.T) {
	u := sampleUnit()
	// Third competitor: no medal, empty mark -> must NOT produce a result.
	u.Competitors = append(
		u.Competitors,
		apiCompetitor{Code: "GBR-003", Noc: "GBR", Name: "C. Slow"},
	)

	out := mapSchedule(apiScheduleResponse{Units: []apiUnit{u}})
	if len(out.results) != 2 {
		t.Fatalf("results = %d, want 2 (empty competitor filtered)", len(out.results))
	}
	wantFiltered := eventcore.NewID(eventcore.ProviderOlympics, "GBR-003")
	for _, r := range out.results {
		if r.ParticipantID == wantFiltered {
			t.Errorf("competitor with no medal/mark must not produce a Result")
		}
	}
}

func TestMapResults_MarkOnlyKept(t *testing.T) {
	u := sampleUnit()
	u.Competitors = u.Competitors[:1]
	u.Competitors[0].Results.MedalType = ""
	u.Competitors[0].Results.Mark = "10.05"
	out := mapResults(u.Competitors, eventcore.NewID(eventcore.ProviderOlympics, "f"))
	if len(out) != 1 {
		t.Fatalf("results = %d, want 1 (mark present)", len(out))
	}
}

func TestMapOutcome(t *testing.T) {
	cases := []struct {
		medal, wlt string
		want       eventcore.Outcome
	}{
		{"GM", "", eventcore.OutcomeMedalGold},
		{"SM", "", eventcore.OutcomeMedalSilver},
		{"BM", "", eventcore.OutcomeMedalBronze},
		{medalCodeGold, "", eventcore.OutcomeMedalGold},
		{medalCodeSilver, "", eventcore.OutcomeMedalSilver},
		{"ME_BRONZE", "", eventcore.OutcomeMedalBronze},
		{"", "W", eventcore.OutcomeWin},
		{"", "L", eventcore.OutcomeLoss},
		{"", "T", eventcore.OutcomeDraw},
		{"", "", eventcore.OutcomeNone},
	}
	for _, tc := range cases {
		if got := mapOutcome(tc.medal, tc.wlt); got != tc.want {
			t.Errorf("mapOutcome(%q,%q) = %q, want %q", tc.medal, tc.wlt, got, tc.want)
		}
	}
}

func TestMapSchedule_ChecksumChangesOnOutcomeChange(t *testing.T) {
	base := mapSchedule(apiScheduleResponse{Units: []apiUnit{sampleUnit()}})
	baseSum := base.fixtures[0].Checksum

	changed := sampleUnit()
	changed.Competitors[0].Results.MedalType = medalCodeSilver // gold -> silver
	out := mapSchedule(apiScheduleResponse{Units: []apiUnit{changed}})

	if out.fixtures[0].Checksum == baseSum {
		t.Errorf("checksum unchanged after outcome change; medal changes must re-notify")
	}
	if baseSum == "" {
		t.Errorf("checksum is empty")
	}
}
