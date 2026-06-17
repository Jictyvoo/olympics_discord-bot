# Imports

## Grouping

Three import groups, separated by blank lines, in this order:

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/wrapped-owls/goremy-di/remy"

    "github.com/jictyvoo/olhojogo/internal/domain/eventcore"
    "github.com/jictyvoo/olhojogo/internal/infra/repositories/reposqlite"
)
```

1. Stdlib.
2. Third-party.
3. This module (`github.com/jictyvoo/olhojogo/...`).

`goimports` enforces this; `task tools:fmt` runs it.

## Dependency direction

```
cmd/<bin>  →  bootstrap  →  domain  ←  infra
                            ▲          ▲
                            └── pkg ───┘
```

- `cmd/<bin>` imports `bootstrap`, `config`, and whichever `infra` packages it needs for
  glue (e.g. `discordgo` session creation in `serve.go`).
- `bootstrap` imports `domain`, `infra/*`, `config`, and the migrator.
- `domain/*` imports other `domain/*` packages and `pkg/*`. Never `infra`.
- `infra/*` imports `domain/*` to implement ports, and `pkg/*`. Never `bootstrap` or `cmd`.
- `pkg/*` imports stdlib + third-party. Never anything under `internal/`.

`go vet ./...` flags any reverse arrow. CI fails on them.

## Blank imports

Allowed only:

- `cmd/<binary>/main.go` — SQL drivers (`_ "modernc.org/sqlite"`, `_ "github.com/go-sql-driver/mysql"`).
- `tools/tools.go` — build-time tool binaries under `//go:build tools`.

Nowhere else.

## Aliasing

- Don't alias unless two packages have the same final path segment (then alias the less canonical).
- `appconfig "github.com/jictyvoo/olhojogo/config"` only inside files that also import a
  `config` from another module.

## Cyclic imports

If you hit a cycle, the design is wrong. Common fixes:

- Move the shared type into a smaller, leaf package (often `eventcore` or a new sub-package).
- Declare the interface in the consumer, satisfied structurally by the provider.
- Never use blank-import tricks or `init()` registration to break cycles — fix the layout.

## What NOT to do

- No `import . "foo"` dot imports.
- No relative imports.
- No `vendor/`. Module mode only.
