package fifafetch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

const (
	defaultBaseURL = "https://api.fifa.com/api/v3"
	defaultLang    = "en"
)

type Provider struct {
	client        httpdatasource.Client
	cache         cachestore.Cache
	baseURL       string
	lang          string
	competitionID string
	seasonID      string
}

func New(
	client httpdatasource.Client,
	cache cachestore.Cache,
	baseURL, lang, competitionID, seasonID string,
) Provider {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if lang == "" {
		lang = defaultLang
	}
	return Provider{
		client:        client,
		cache:         cache,
		baseURL:       baseURL,
		lang:          lang,
		competitionID: competitionID,
		seasonID:      seasonID,
	}
}

func (p Provider) Code() eventcore.ProviderID { return eventcore.ProviderFIFA }
func (p Provider) DisplayName() string        { return "FIFA" }

func (p Provider) SyncFixturesByDate(
	ctx context.Context,
	day time.Time,
) (eventcore.SyncDelta, error) {
	url := matchesByDayURL(p.baseURL, p.competitionID, p.seasonID, p.lang, day)
	body, err := p.fetch(ctx, url)
	if err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("fifa: fetch matches: %w", err)
	}

	var resp apiMatchesResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("fifa: decode matches: %w", err)
	}

	season := p.fetchSeason(ctx, day)
	mapped := mapMatches(resp, p.lang, season)
	standings, stErr := p.fetchStandings(ctx, mapped.stageKeys)

	delta := eventcore.SyncDelta{
		Competitions: mapped.competitions,
		Seasons:      mapped.seasons,
		Stages:       mapped.stages,
		Groups:       mapped.groups,
		Venues:       mapped.venues,
		Participants: mapped.participants,
		Fixtures:     mapped.fixtures,
		Results:      mapped.results,
		Standings:    standings,
	}
	return delta, stErr
}

func (p Provider) SyncFixtureResults(
	_ context.Context,
	_ eventcore.Fixture,
) (eventcore.SyncDelta, error) {
	return eventcore.SyncDelta{}, eventcore.ErrNotImplemented
}
