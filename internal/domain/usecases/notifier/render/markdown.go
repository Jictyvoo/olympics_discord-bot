package render

import (
	"fmt"
	"strings"
)

// Markdown renders a fixture as a Discord-compatible markdown block.
type Markdown struct{}

func (Markdown) Render(view FixtureView) string {
	fixture := view.Fixture
	var b strings.Builder
	fmt.Fprintf(&b, "**%s**\n", fixture.Name)
	fmt.Fprintf(&b, "Status: %s\n", fixture.Status)
	fmt.Fprintf(&b, "Start: %s\n", DiscordTimestamp(fixture.StartsAt))
	fmt.Fprintf(&b, "End: %s\n", DiscordTimestamp(fixture.EndsAt))
	if len(view.Competitors) > 0 {
		b.WriteString("Participants:")
		for _, c := range view.Competitors {
			fmt.Fprintf(&b, "\n- %s", c.Participant.Name)
		}
	} else if len(fixture.Participants) > 0 {
		b.WriteString("Participants:")
		for _, p := range fixture.Participants {
			fmt.Fprintf(&b, "\n- `%s`", p.ParticipantID.String()[:8])
		}
	}
	return b.String()
}
