package bootstrap

import (
	"log/slog"

	"github.com/wrapped-owls/goremy-di/remy"

	appconfig "github.com/jictyvoo/olhojogo/config"
	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/domain/provider"
	"github.com/jictyvoo/olhojogo/internal/infra/fetchers/fifafetch"
	"github.com/jictyvoo/olhojogo/internal/infra/fetchers/olympicsfetch"
)

func registerProviders(inj remy.Injector, conf appconfig.Config) {
	enabled := make([]eventcore.ProviderID, 0, len(conf.Providers))
	for _, pc := range conf.Providers {
		if !pc.Enabled {
			continue
		}
		switch pc.Code {
		case eventcore.ProviderOlympics:
			olympicsfetch.Register(inj, pc.BaseURL, pc.Language)
		case eventcore.ProviderFIFA:
			fifafetch.Register(inj, pc.BaseURL, pc.Language, pc.CompetitionID, pc.SeasonID)
		default:
			slog.Warn("bootstrap: unknown provider code", slog.String("code", pc.Code))
			continue
		}
		enabled = append(enabled, pc.Code)
	}
	provider.RegisterSet(inj, enabled)
}
