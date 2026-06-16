package bootstrap

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore"
)

func TestRegisterCache(t *testing.T) {
	tests := []struct {
		name    string
		backend string
	}{
		{name: "memory backend", backend: "memory"},
		{name: "file backend", backend: "file"},
		{name: "default backend", backend: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := appconfig.Config{}
			conf.Cache.Backend = tt.backend
			conf.Cache.FilePath = t.TempDir()

			inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
			registerCache(inj, conf)

			c, err := remy.Get[cachestore.Cache](inj)
			if err != nil {
				t.Fatalf("resolving Cache failed: %v", err)
			}
			if c == nil {
				t.Fatal("resolved Cache is nil")
			}
		})
	}
}
