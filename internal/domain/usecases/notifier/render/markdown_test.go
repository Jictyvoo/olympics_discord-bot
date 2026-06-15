package render

import (
	"strings"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func TestMarkdownRender(t *testing.T) {
	start := time.Date(2024, 7, 26, 20, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	tests := []struct {
		name        string
		view        FixtureView
		wantSubstrs []string
	}{
		{
			name: "basic fixture with named participants",
			view: FixtureView{
				Fixture: eventcore.Fixture{
					Name:     "100m Final",
					Status:   eventcore.FixtureScheduled,
					StartsAt: start,
					EndsAt:   end,
				},
				Participants: []eventcore.Participant{
					{Name: "Usain Bolt"},
					{Name: "Carl Lewis"},
				},
			},
			wantSubstrs: []string{
				"**100m Final**",
				"Status: scheduled",
				"Start: <t:1722024000:R>",
				"End: <t:1722031200:R>",
				"- Usain Bolt",
				"- Carl Lewis",
			},
		},
		{
			name: "uses relative discord timestamps not raw unix",
			view: FixtureView{
				Fixture: eventcore.Fixture{Name: "Match", StartsAt: start, EndsAt: end},
			},
			wantSubstrs: []string{"<t:1722024000:R>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Markdown{}.Render(tt.view)
			for _, want := range tt.wantSubstrs {
				if !strings.Contains(got, want) {
					t.Fatalf("output missing %q\ngot:\n%s", want, got)
				}
			}
		})
	}
}
