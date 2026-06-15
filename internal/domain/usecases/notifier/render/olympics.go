package render

import (
	"fmt"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/pkg/strutil"
)

// Olympics extends the default Markdown renderer with medal- and
// discipline-aware output derived from the render aggregate.
type Olympics struct{ Markdown }

func (Olympics) Render(view FixtureView) string {
	fixture := view.Fixture
	competition := view.Competition

	var b strings.Builder

	discipline := competition.Discipline
	if discipline == "" {
		discipline = competition.Name
	}
	if icon := DisciplineIcon(competition.Code); icon != "" {
		fmt.Fprintf(&b, "# %s %s\n", icon, discipline)
	} else {
		fmt.Fprintf(&b, "# %s\n", discipline)
	}

	medalIcon := ""
	if hasMedal(view.Results) {
		medalIcon = ":medal: "
	}
	fmt.Fprintf(&b, "**Event:** %s%s", medalIcon, fixture.Name)
	if fixture.Status != "" {
		fmt.Fprintf(&b, " - %s", fixture.Status)
	}
	b.WriteByte('\n')

	// Phase: suppress when it merely repeats the event name.
	if phase := discipline; phase != "" && !strutil.EqualAlfaNum(fixture.Name, phase) {
		fmt.Fprintf(&b, "**Phase:** %s\n", phase)
	}

	fmt.Fprintf(&b, "**Start:** %s\n", DiscordTimestamp(fixture.StartsAt))
	fmt.Fprintf(&b, "**End:** %s\n", DiscordTimestamp(fixture.EndsAt))

	if len(view.Participants) > 0 {
		b.WriteString("**Competitors:**\n")
		for _, p := range view.Participants {
			fmt.Fprintf(&b, "- %s\n", p.Name)
		}
	}

	return b.String()
}

// hasMedal reports whether any result carries a medal outcome.
func hasMedal(results []eventcore.Result) bool {
	for _, r := range results {
		switch r.Outcome {
		case eventcore.OutcomeMedalGold,
			eventcore.OutcomeMedalSilver,
			eventcore.OutcomeMedalBronze:
			return true
		}
	}
	return false
}
