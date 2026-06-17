# Patterns

Implementation recipes with copyable templates. These are not mandatory in the way `rules/` are
— a pattern is a known-good recipe for a recurring problem.

| File                                             | Topic                                              |
|--------------------------------------------------|----------------------------------------------------|
| [provider-onboarding.md](provider-onboarding.md) | Adding a new HTTP fetcher (Strategy impl)          |
| [repo-layout.md](repo-layout.md)                 | sqlc + repo struct + mapper recipe per dialect     |
| [bootstrap-di.md](bootstrap-di.md)               | The composition root and dialect/fetcher switching |
| [config-layout.md](config-layout.md)             | Typed Config + env binder pattern                  |
| [observer-weakptr.md](observer-weakptr.md)       | weak.Pointer observer registry                     |
| [discord-sync.md](discord-sync.md)               | Discord Scheduled Events lifecycle                 |
| [mocks.md](mocks.md)                             | Colocated `_mock_test.go` recipe                   |
| [migrations.md](migrations.md)                   | Embedded SQL migrations + in-house runner          |
