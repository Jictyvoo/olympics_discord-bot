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
	competition := view.Context.Competition

	var b strings.Builder

	discipline := competition.Discipline
	if discipline == "" {
		discipline = competition.Name
	}
	icon := DisciplineIcon(competition.Code)
	if icon == "" {
		icon = DisciplineIconByName(discipline)
	}
	if icon != "" {
		fmt.Fprintf(&b, "# %s %s\n", icon, discipline)
	} else {
		fmt.Fprintf(&b, "# %s\n", discipline)
	}

	medalIcon := ""
	if hasMedal(view.Competitors) {
		medalIcon = ":medal: "
	}
	fmt.Fprintf(&b, "**Event:** %s%s", medalIcon, fixture.Name)
	if fixture.Status != "" {
		fmt.Fprintf(&b, " - %s", fixture.Status)
	}
	b.WriteByte('\n')

	// Suppress the phase when it merely repeats the event name.
	if phase := phaseLabel(view.Context); phase != "" &&
		!strutil.EqualAlfaNum(fixture.Name, phase) {
		fmt.Fprintf(&b, "**Phase:** %s\n", phase)
	}

	if gender := genderLabel(view.Competitors); gender != "" {
		fmt.Fprintf(&b, "**Gender:** %s\n", gender)
	}

	fmt.Fprintf(&b, "**Start:** %s\n", DiscordTimestamp(fixture.StartsAt))
	fmt.Fprintf(&b, "**End:** %s\n", DiscordTimestamp(fixture.EndsAt))

	if len(view.Competitors) > 0 {
		b.WriteString("**Competitors:**\n")
		for _, c := range view.Competitors {
			b.WriteString(competitorLine(c))
		}
	}

	return b.String()
}

func competitorLine(c eventcore.FixtureCompetitor) string {
	var b strings.Builder
	b.WriteString("- ")
	if flag := (eventcore.Country{ISO2: c.CountryISO2}).EmojiFlag(); flag != "" {
		b.WriteString(flag)
		b.WriteByte(' ')
	}
	b.WriteString(c.Participant.Name)
	if c.Result.Score != "" {
		fmt.Fprintf(&b, " #%s", c.Result.Score)
	}
	if mark := outcomeIcon(c.Result.Outcome); mark != "" {
		fmt.Fprintf(&b, " (%s)", mark)
	}
	b.WriteByte('\n')
	return b.String()
}

func phaseLabel(ctx eventcore.FixtureContext) string {
	parts := make([]string, 0, 2) //nolint:mnd // stage + group
	if ctx.StageName != "" {
		parts = append(parts, ctx.StageName)
	}
	if ctx.GroupName != "" {
		parts = append(parts, ctx.GroupName)
	}
	return strings.Join(parts, " - ")
}

func genderLabel(competitors []eventcore.FixtureCompetitor) string {
	for _, c := range competitors {
		switch c.Participant.Gender {
		case "M":
			return "Male"
		case "F":
			return "Female"
		case "mixed":
			return "Mixed"
		}
	}
	return ""
}

// outcomeIcon reads a win or gold as first place, a loss or silver as second.
func outcomeIcon(o eventcore.Outcome) string {
	switch o {
	case eventcore.OutcomeMedalGold, eventcore.OutcomeWin:
		return ":first_place:"
	case eventcore.OutcomeMedalSilver, eventcore.OutcomeLoss:
		return ":second_place:"
	case eventcore.OutcomeMedalBronze:
		return ":third_place:"
	case eventcore.OutcomeDraw:
		return ":handshake:"
	}
	return ""
}

func hasMedal(competitors []eventcore.FixtureCompetitor) bool {
	for _, c := range competitors {
		switch c.Result.Outcome {
		case eventcore.OutcomeMedalGold,
			eventcore.OutcomeMedalSilver,
			eventcore.OutcomeMedalBronze:
			return true
		}
	}
	return false
}
