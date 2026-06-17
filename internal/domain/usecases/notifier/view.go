package notifier

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier/render"
)

// compose appends the @mention line after the body so the ping reads as a footer.
func (n *Notifier) compose(f eventcore.Fixture) (string, error) {
	view := n.buildView(f)
	content := n.renderer.Render(view)
	suffix, err := n.mentionSuffix(view)
	if err != nil {
		return "", err
	}
	if suffix != "" {
		content = strings.TrimRight(content, "\n") + "\n" + suffix
	}
	return content, nil
}

// buildView assembles the render aggregate for a fixture. Enrichment reads are
// best-effort: a failure is logged and the view degrades gracefully rather than
// blocking the notification, whose essential payload is the fixture itself.
func (n *Notifier) buildView(f eventcore.Fixture) render.FixtureView {
	view := render.FixtureView{Fixture: f}

	context, err := n.context.GetFixtureContext(f.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Warn(
			"notifier: load context",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	} else {
		view.Context = context
	}

	competitors, err := n.competitors.ListFixtureCompetitors(f.ID)
	if err != nil {
		slog.Warn(
			"notifier: load competitors",
			slog.String("fixture", f.ID.String()),
			slog.String("err", err.Error()),
		)
	} else {
		view.Competitors = competitors
	}

	return view
}

// mentionSuffix renders the matching subscribers as a "<@id> <@id>" line,
// or "" when there are none.
func (n *Notifier) mentionSuffix(
	view render.FixtureView,
) (string, error) {
	if n.mentions == nil {
		return "", nil
	}
	countryCodes := make([]string, 0, len(view.Competitors))
	for _, c := range view.Competitors {
		if c.Participant.CountryISO != "" {
			countryCodes = append(countryCodes, c.Participant.CountryISO)
		}
	}
	users, err := n.mentions.MentionsFor(n.guildID, countryCodes, view.Context.Competition.Code)
	if err != nil {
		return "", err
	}
	if len(users) == 0 {
		return "", nil
	}
	parts := make([]string, len(users))
	for i, id := range users {
		parts[i] = fmt.Sprintf("<@%s>", id)
	}
	return strings.Join(parts, " "), nil
}
