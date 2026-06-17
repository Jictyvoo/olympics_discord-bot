# Commits

Format: `:emoji-name: type(scope): Subject`.

- `:emoji-name:` — gitmoji label, lowercase, surrounded by colons.
- `type` — conventional commit type: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`, `style`,
  `perf`, `build`, `ci`.
- `(scope)` — optional; the package or area being changed.
- `Subject` — uppercase imperative, no trailing period, ≤90 chars total line length.

## Examples

```
:sparkles: feat(provider): Add FIFA fetcher
:bug: fix(notifier): Honour notify window when checksum changes
:recycle: refactor(repo): Route DI through consumer interfaces
:wrench: chore(taskfile): Split gen targets per dialect
:page_with_curl: docs(rules): Rewrite for olhojogo layout
:test_tube: test(syncer): Cover SyncDay edge cases
```

## Body

Optional. When present:

- Wrap at 72 chars.
- Explain the *why*, not the *what*.
- Reference prior work or related issues only when they constrain the choice.
- No co-author trailers unless the user requests one.

## Picking emoji

Vary the emoji to reflect the change. Don't reuse `:recycle:` for every refactor commit — match
the work. Reference:

| Change                         | Emoji                         |
|--------------------------------|-------------------------------|
| New user-visible feature       | `:sparkles:`                  |
| Bug fix                        | `:bug:`                       |
| Refactor (no behaviour change) | `:recycle:`                   |
| Docs                           | `:page_with_curl:`            |
| Tests                          | `:test_tube:`                 |
| Tooling / config               | `:wrench:`                    |
| Delete code                    | `:fire:`                      |
| Style / formatting             | `:art:`                       |
| Performance                    | `:zap:`                       |
| Move/rename files              | `:truck:`                     |
| New module / package           | `:package:`                   |
| Dependencies                   | `:arrow_up:` / `:arrow_down:` |
| Cleanup                        | `:broom:`                     |

## What NOT to do

- No `chore: update` with no context.
- No "wip" commits in shared branches.
- No commits that mix a refactor with a behaviour change. Split them.
- No `--no-verify` to skip hooks.
- No `Co-Authored-By: Claude` (or any AI co-author) unless explicitly requested.
