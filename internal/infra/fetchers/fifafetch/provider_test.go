package fifafetch

import (
	"errors"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
)

func TestProvider_Identity(t *testing.T) {
	p := New()
	if p.Code() != eventcore.ProviderFIFA {
		t.Errorf("Code = %q, want %q", p.Code(), eventcore.ProviderFIFA)
	}
	if p.DisplayName() == "" {
		t.Error("DisplayName must not be empty")
	}
}

func TestProvider_StubReturnsNotImplemented(t *testing.T) {
	p := New()
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); !errors.Is(
		err, eventcore.ErrNotImplemented,
	) {
		t.Errorf("SyncFixturesByDate: want ErrNotImplemented, got %v", err)
	}
	if _, err := p.SyncFixtureResults(t.Context(), eventcore.Fixture{}); !errors.Is(
		err, eventcore.ErrNotImplemented,
	) {
		t.Errorf("SyncFixtureResults: want ErrNotImplemented, got %v", err)
	}
}
