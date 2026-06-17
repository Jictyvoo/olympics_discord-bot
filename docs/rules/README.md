# Rules

Mandatory rules for all new and modified code in this repository.
If a rule conflicts with code you find, the code is wrong — update it.

| File                                   | Topic                                                                                                     |
|----------------------------------------|-----------------------------------------------------------------------------------------------------------|
| [code-placement.md](code-placement.md) | What goes in `cmd/`, `config/`, `internal/{bootstrap,domain,infra,migrator}/`, `pkg/`, `build/`, `tools/` |
| [naming.md](naming.md)                 | Package names, type names, file names                                                                     |
| [effective-go.md](effective-go.md)     | The Go style subset we enforce                                                                            |
| [types.md](types.md)                   | Domain types, IDs, value vs pointer receivers                                                             |
| [imports.md](imports.md)               | Import grouping and dependency direction                                                                  |
| [errors.md](errors.md)                 | `errors.Join`, `%w`, sentinel pattern                                                                     |
| [logging.md](logging.md)               | `log/slog` only; no wrappers                                                                              |
| [config.md](config.md)                 | Typed config + env binders; `os.Getenv` ban                                                               |
| [startup.md](startup.md)               | No I/O at `init()`; main is the only entrypoint                                                           |
| [wiring.md](wiring.md)                 | remy DI: `RegisterConstructor*` only; bootstrap-driven switches                                           |
| [http-clients.md](http-clients.md)     | Outbound HTTP through `httpdatasource` only                                                               |
| [concurrency.md](concurrency.md)       | Goroutine ownership and cancellation                                                                      |
| [testing.md](testing.md)               | Table-driven, colocated, mocks pattern                                                                    |
| [security.md](security.md)             | Secrets, TLS, untrusted input                                                                             |
| [commits.md](commits.md)               | Gitmoji + Conventional Commits                                                                            |
| [build.md](build.md)                   | Taskfile entry points; CGO-free default                                                                   |

See also: [`../patterns/`](../patterns/) for implementation recipes.
