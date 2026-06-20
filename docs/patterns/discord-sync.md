# Discord Scheduled Events — lifecycle

`internal/domain/usecases/discordsync` manages the create/update/cancel lifecycle of Discord
Scheduled Events for fixtures within the configured `discord.horizon` window (default 14 days).

## State machine

```
fixture.status            discord_events row           Discord Scheduled Event
─────────────────────────────────────────────────────────────────────────────
scheduled / live          missing                  →   Create (write back discord_event_id)
scheduled / live          exists, checksum changed →   Update (title, description, time)
cancelled / postponed     exists                   →   Cancel
finished                  exists                   →   no-op (already completed in Discord)
```

## Idempotency

- `(fixture_id, guild_id)` is the PK on `discord_events`.
- `UNIQUE(discord_event_id)` prevents duplicate creates.
- `last_checksum` column tracks the last-seen fixture checksum; an update fires only when it
  changes.

## Trigger

`DiscordSync.Run(ctx)` is called by the observer hook after each successful provider sync. It
has no ticker of its own — it piggybacks on the sync interval.

```go
func (d *DiscordSync) Run(ctx context.Context) error {
    fixtures, err := d.fixtures.ListFixturesStartingBefore(ctx, time.Now().Add(d.horizon))
    if err != nil {
        return err
    }
    for _, f := range fixtures {
        if err := d.reconcile(ctx, f); err != nil {
            slog.Error("discordsync: reconcile", slog.String("fixture_id", f.ID.String()), slog.String("err", err.Error()))
        }
    }
    return nil
}
```

## ScheduledEventInput

```go
type ScheduledEventInput struct {
    Name        string
    Description string
    StartsAt    time.Time
    EndsAt      time.Time
    ChannelID   string  // voice channel the event is associated with
    Location    string  // for external-location events
}
```

`builder.go` constructs one per fixture and resolves channel IDs from config via
`discordfacade.Facade.ResolveChannel`.

## Adding a guild

- `config.Discord.GuildID` and `DefaultChannel` are required.
- Scheduled events are guild-wide, not channel-bound, so they ignore per-provider channels.
  The `config.ProviderCfg.DiscordChannel` override only routes the notifier's text messages: a
  fixture goes to its provider's channel, falling back to `DefaultChannel` (see `notifier.channelRouter`).
- Channel names resolve to IDs at startup via `discordfacade.Facade.ResolveChannel`, cached in
  the facade. Never hardcode channel IDs or names in use-case code.
