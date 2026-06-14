package cachestore

import (
	"time"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/infra/cachestore/filecache"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore/memcache"
)

func RegisterMemory(inj remy.Injector, ttl time.Duration) {
	remy.RegisterConstructor(inj, remy.Factory[Cache],
		func() Cache { return memcache.New(ttl) })
}

func RegisterFile(inj remy.Injector, rootPath string, ttl time.Duration) {
	remy.RegisterConstructorErr(inj, remy.Factory[Cache],
		func() (Cache, error) { return filecache.New(rootPath, ttl) })
}
