package repoinj

import (
	"reflect"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

type Location struct {
	Code string
}

type Battlegear interface {
	LocationCode() string
}

type CastleBodhran struct {
	location Location
}

func (c *CastleBodhran) LocationCode() string {
	return c.location.Code
}

type MugicCounter struct{}

const castleBodhranCode = "castle-bodhran"

//nolint:ireturn // factory returning consumer interface by design
func newContainer() remy.Injector {
	container := remy.NewInjector()
	remy.RegisterInstance(
		container,
		Location{
			Code: castleBodhranCode,
		},
	)
	return container
}

func TestRegisterAliased(t *testing.T) {
	tests := []struct {
		name    string
		kind    InjectionKind
		wantErr bool
	}{
		{
			name: "factory",
			kind: InjectionFactory,
		},
		{
			name: "singleton",
			kind: InjectionSingleton,
		},
		{
			name: "lazy singleton",
			kind: InjectionLazySingleton,
		},
		{
			name:    "invalid kind",
			kind:    InjectionKind(255),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runRegisterAliasedCase(t, tt.kind, tt.wantErr)
		})
	}
}

func runRegisterAliasedCase(t *testing.T, kind InjectionKind, wantErr bool) {
	t.Helper()
	container := newContainer()

	err := RegisterAliased[*CastleBodhran, Location, Battlegear](
		container,
		kind,
		func(location Location) *CastleBodhran {
			return &CastleBodhran{location: location}
		},
	)

	if wantErr {
		if err == nil {
			t.Fatal("expected error")
		}
		return
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertResolvedAlias(t, container)
}

func assertResolvedAlias(t *testing.T, container remy.Injector) {
	t.Helper()

	concrete, err := remy.Get[*CastleBodhran](container)
	if err != nil {
		t.Fatalf("failed resolving concrete type: %v", err)
	}

	alias, err := remy.Get[Battlegear](container)
	if err != nil {
		t.Fatalf("failed resolving alias type: %v", err)
	}

	if concrete.LocationCode() != castleBodhranCode {
		t.Fatalf(
			"expected injected location code %s, got %q",
			castleBodhranCode,
			concrete.LocationCode(),
		)
	}

	if alias.LocationCode() != castleBodhranCode {
		t.Fatalf(
			"expected injected location code %s, got %q",
			castleBodhranCode,
			alias.LocationCode(),
		)
	}

	if !reflect.DeepEqual(concrete, alias) {
		t.Fatalf(
			"expected concrete and alias resolutions to be equivalent\nconcrete=%#v\nalias=%#v",
			concrete,
			alias,
		)
	}
}

func TestRegisterAliased_TypeDoesNotImplementInterface(t *testing.T) {
	container := newContainer()

	err := RegisterAliased[MugicCounter, Location, Battlegear](
		container,
		InjectionFactory,
		func(Location) MugicCounter {
			return MugicCounter{}
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
