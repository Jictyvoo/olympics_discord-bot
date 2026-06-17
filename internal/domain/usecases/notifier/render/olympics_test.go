package render

import (
	"strings"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

const (
	athleticsCode = "ATH"
	athleticsName = "Athletics"
	futebol       = "Futebol"
	finalName     = "Final"
)

type olympicsRenderCase struct {
	name        string
	view        FixtureView
	wantSubstrs []string
	notSubstrs  []string
}

func assertOlympicsRender(t *testing.T, tt olympicsRenderCase) {
	t.Helper()
	got := Olympics{}.Render(tt.view)
	for _, want := range tt.wantSubstrs {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q\ngot:\n%s", want, got)
		}
	}
	for _, bad := range tt.notSubstrs {
		if strings.Contains(got, bad) {
			t.Fatalf("output should not contain %q\ngot:\n%s", bad, got)
		}
	}
}

//nolint:funlen // table-driven fixture
func olympicsRenderCases(start, end time.Time) []olympicsRenderCase {
	return []olympicsRenderCase{
		{
			name: "medal prefix derived from gold result",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: finalName, StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{
						Code:       athleticsCode,
						Discipline: athleticsName,
					},
				},
				Competitors: []eventcore.FixtureCompetitor{
					{Result: eventcore.Result{Outcome: eventcore.OutcomeMedalGold}},
				},
			},
			wantSubstrs: []string{":medal:", athleticsIcon, athleticsName},
		},
		{
			name: "no medal prefix for non-medal outcomes",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: "Heat 1", StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{Code: "SWM", Discipline: "Swimming"},
				},
				Competitors: []eventcore.FixtureCompetitor{
					{Result: eventcore.Result{Outcome: eventcore.OutcomeWin}},
				},
			},
			wantSubstrs: []string{":person_swimming:", "Swimming"},
			notSubstrs:  []string{":medal:"},
		},
		{
			name: "unknown discipline code shows name without icon",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: "Event", StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{Code: "ZZZ", Discipline: "Mystery Sport"},
				},
			},
			wantSubstrs: []string{"# Mystery Sport"},
		},
		{
			name: "football icon resolved by discipline name when code is unknown",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: finalName, StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{Code: "17", Discipline: futebol},
				},
			},
			wantSubstrs: []string{"# :soccer: Futebol"},
		},
		{
			name: "redundant phase suppressed via alfa-num equality",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: athleticsName, StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{
						Code:       athleticsCode,
						Discipline: athleticsName,
					},
					StageName: "Athletics!",
				},
			},
			notSubstrs: []string{"**Phase:**"},
		},
		{
			name: "phase joins stage and group",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: "Match", StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{Code: "17", Discipline: futebol},
					StageName:   "Primeira fase",
					GroupName:   "Group A",
				},
			},
			wantSubstrs: []string{"**Phase:** Primeira fase - Group A"},
		},
		{
			name: "competitors render flag, score, placement and gender",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: finalName, StartsAt: start, EndsAt: end},
				Context: eventcore.FixtureContext{
					Competition: eventcore.Competition{Code: "17", Discipline: futebol},
				},
				Competitors: []eventcore.FixtureCompetitor{
					{
						Participant: eventcore.Participant{Name: "EUA", Gender: "F"},
						CountryISO2: "US",
						Result:      eventcore.Result{Score: "1", Outcome: eventcore.OutcomeWin},
					},
					{
						Participant: eventcore.Participant{Name: "Brasil", Gender: "F"},
						CountryISO2: "BR",
						Result:      eventcore.Result{Score: "0", Outcome: eventcore.OutcomeLoss},
					},
				},
			},
			wantSubstrs: []string{
				"**Gender:** Female",
				"**Competitors:**",
				"- :flag_us: EUA #1 (:first_place:)",
				"- :flag_br: Brasil #0 (:second_place:)",
			},
		},
	}
}

func TestOlympicsRender(t *testing.T) {
	start := time.Date(2024, 7, 26, 20, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	for _, tt := range olympicsRenderCases(start, end) {
		t.Run(tt.name, func(t *testing.T) {
			assertOlympicsRender(t, tt)
		})
	}
}
