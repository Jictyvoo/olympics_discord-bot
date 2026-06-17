package render

import "github.com/jictyvoo/olhojogo/internal/domain/eventcore"

// FixtureView is the render aggregate: a fixture plus the related entities a
// renderer needs to produce a notification. The notifier assembles it so that
// Olympics-specific concerns (medals, disciplines) stay out of eventcore.Fixture.
type FixtureView struct {
	Fixture     eventcore.Fixture
	Context     eventcore.FixtureContext
	Competitors []eventcore.FixtureCompetitor
}

// Renderer produces a human-readable notification string for a fixture view.
type Renderer interface {
	Render(view FixtureView) string
}
