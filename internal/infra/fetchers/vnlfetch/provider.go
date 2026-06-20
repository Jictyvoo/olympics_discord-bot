package vnlfetch

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
	defaultBaseURL = "https://br.volleyballworld.com"
	defaultLang    = "en"
)

type Provider struct {
	client      httpdatasource.Client
	cache       cachestore.Cache
	baseURL     string
	lang        string
	tournaments string // semicolon-separated tournament numbers, e.g. "1661;1662"
}

func New(
	client httpdatasource.Client,
	cache cachestore.Cache,
	baseURL, lang, tournaments string,
) Provider {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if lang == "" {
		lang = defaultLang
	}
	return Provider{
		client:      client,
		cache:       cache,
		baseURL:     baseURL,
		lang:        lang,
		tournaments: tournaments,
	}
}

func (p Provider) Code() eventcore.ProviderID { return eventcore.ProviderVNL }
func (p Provider) DisplayName() string        { return "Volleyball Nations League" }

func (p Provider) SyncFixturesByDate(
	ctx context.Context,
	day time.Time,
) (eventcore.SyncDelta, error) {
	url := scheduleByDayURL(p.baseURL, p.tournaments, day)
	body, err := p.fetch(ctx, url)
	if err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("vnl: fetch schedule: %w", err)
	}

	var resp apiScheduleResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("vnl: decode schedule: %w", err)
	}

	mapped := mapSchedule(resp, p.lang, day, time.Now().UTC())
	delta := eventcore.SyncDelta{
		Competitions: mapped.competitions,
		Seasons:      mapped.seasons,
		Stages:       mapped.stages,
		Groups:       mapped.groups,
		Venues:       mapped.venues,
		Participants: mapped.participants,
		Fixtures:     mapped.fixtures,
		Results:      mapped.results,
	}
	return delta, nil
}

func (p Provider) SyncFixtureResults(
	_ context.Context,
	_ eventcore.Fixture,
) (eventcore.SyncDelta, error) {
	return eventcore.SyncDelta{}, eventcore.ErrNotImplemented
}
