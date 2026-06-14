package reposqlite

import (
	"context"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestRepoCtx_UsesInjectedContext(t *testing.T) {
	type ctxKey struct{}
	parent := context.WithValue(context.Background(), ctxKey{}, "marker")

	base := newRepo(parent, nil)
	qctx, cancel := base.Ctx()
	defer cancel()

	if got := qctx.Value(ctxKey{}); got != "marker" {
		t.Fatalf("Ctx did not derive from injected context: value=%v", got)
	}
	if _, ok := qctx.Deadline(); !ok {
		t.Fatal("Ctx must carry the default query timeout deadline")
	}
}

func TestRepoCtx_NilFallsBackToBackground(t *testing.T) {
	base := newRepo(nil, nil) //nolint:staticcheck // intentionally exercising the nil-ctx fallback
	qctx, cancel := base.Ctx()
	defer cancel()
	if qctx == nil {
		t.Fatal("Ctx returned a nil context")
	}
	if _, ok := qctx.Deadline(); !ok {
		t.Fatal("fallback Ctx must still carry a timeout deadline")
	}
}

// TestGetWithContext_FlowsCtxToBase: GetWithContext flows the per-call ctx into
// the base, while plain Get falls back to the registered Background instance.
func TestGetWithContext_FlowsCtxToBase(t *testing.T) {
	type ctxKey struct{}
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance[context.Context](inj, context.Background())
	Register(inj, nil)

	// Plain Get: resolves via the fallback Background context (no marker).
	repoFallback, err := remy.Get[FixtureRepo](inj)
	if err != nil {
		t.Fatalf("plain Get: %v", err)
	}
	if v := repoFallback.Context().Value(ctxKey{}); v != nil {
		t.Fatalf("plain Get must use fallback ctx, got marker=%v", v)
	}

	// GetWithContext: the per-call ctx must reach the base.
	want := context.WithValue(context.Background(), ctxKey{}, "tick")
	repoScoped, err := remy.GetWithContext[FixtureRepo](inj, want)
	if err != nil {
		t.Fatalf("GetWithContext: %v", err)
	}
	if v := repoScoped.Context().Value(ctxKey{}); v != "tick" {
		t.Fatalf("GetWithContext did not flow ctx to base: marker=%v", v)
	}

	qctx, cancel := repoScoped.Ctx()
	defer cancel()
	if _, ok := qctx.Deadline(); !ok {
		t.Fatal("scoped Ctx must carry a timeout deadline")
	}
}
