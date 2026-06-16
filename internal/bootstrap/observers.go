package bootstrap

import (
	"log/slog"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/services"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/discordsync"
	"github.com/jictyvoo/olhojogo/internal/domain/usecases/notifier"
)

// WireObservers subscribes the Notifier and DiscordSync to the syncer's fixture
// events. The injector holds strong refs so the subject's weak pointers don't
// evict them.
func WireObservers(inj remy.Injector) {
	subject, err := remy.Get[*services.Subject[eventcore.Fixture]](inj)
	if err != nil {
		slog.Error("bootstrap: get events subject", slog.String("err", err.Error()))
		return
	}

	fixtureNotifier, err := remy.Get[*notifier.Notifier](inj)
	if err != nil {
		slog.Error("bootstrap: get notifier", slog.String("err", err.Error()))
		return
	}
	discordSync, err := remy.Get[*discordsync.DiscordSync](inj)
	if err != nil {
		slog.Error("bootstrap: get discordsync", slog.String("err", err.Error()))
		return
	}

	subject.Register(fixtureNotifier)
	subject.Register(discordSync)
}
