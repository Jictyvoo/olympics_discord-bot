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

func olympicsRenderCases(start, end time.Time) []olympicsRenderCase {
	return []olympicsRenderCase{
		{
			name: "medal prefix derived from gold result",
			view: FixtureView{
				Fixture:     eventcore.Fixture{Name: "Final", StartsAt: start, EndsAt: end},
				Competition: eventcore.Competition{Code: athleticsCode, Discipline: athleticsName},
				Results: []eventcore.Result{
					{Outcome: eventcore.OutcomeMedalGold},
				},
			},
			wantSubstrs: []string{":medal:", athleticsIcon, athleticsName},
		},
		{
			name: "no medal prefix for non-medal outcomes",
			view: FixtureView{
				Fixture:     eventcore.Fixture{Name: "Heat 1", StartsAt: start, EndsAt: end},
				Competition: eventcore.Competition{Code: "SWM", Discipline: "Swimming"},
				Results: []eventcore.Result{
					{Outcome: eventcore.OutcomeWin},
				},
			},
			wantSubstrs: []string{":person_swimming:", "Swimming"},
			notSubstrs:  []string{":medal:"},
		},
		{
			name: "unknown discipline code shows name without icon",
			view: FixtureView{
				Fixture:     eventcore.Fixture{Name: "Event", StartsAt: start, EndsAt: end},
				Competition: eventcore.Competition{Code: "ZZZ", Discipline: "Mystery Sport"},
			},
			wantSubstrs: []string{"# Mystery Sport"},
		},
		{
			name: "redundant phase suppressed via alfa-num equality",
			view: FixtureView{
				Fixture:     eventcore.Fixture{Name: athleticsName, StartsAt: start, EndsAt: end},
				Competition: eventcore.Competition{Code: athleticsCode, Discipline: "Athletics!"},
			},
			notSubstrs: []string{"**Phase:**"},
		},
		{
			name: "competitors listed by name",
			view: FixtureView{
				Fixture:     eventcore.Fixture{Name: "Final", StartsAt: start, EndsAt: end},
				Competition: eventcore.Competition{Code: "BOX", Discipline: "Boxing"},
				Participants: []eventcore.Participant{
					{Name: "Fighter A"},
					{Name: "Fighter B"},
				},
			},
			wantSubstrs: []string{"**Competitors:**", "- Fighter A", "- Fighter B"},
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
