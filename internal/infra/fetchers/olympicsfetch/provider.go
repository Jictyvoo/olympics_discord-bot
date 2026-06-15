package olympicsfetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/cachestore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

const defaultBaseURL = "https://sph-s-api.olympicsfetch.com"

type Provider struct {
	client  httpdatasource.Client
	cache   cachestore.Cache
	baseURL string
	lang    string
}

func New(client httpdatasource.Client, cache cachestore.Cache, baseURL, lang string) Provider {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if lang == "" {
		lang = "ENG"
	}
	return Provider{client: client, cache: cache, baseURL: baseURL, lang: lang}
}

func (p Provider) Code() eventcore.ProviderID { return eventcore.ProviderOlympics }
func (p Provider) DisplayName() string        { return "Olympics" }

func (p Provider) SyncFixturesByDate(
	ctx context.Context,
	day time.Time,
) (eventcore.SyncDelta, error) {
	url := scheduleByDayURL(p.baseURL, p.lang, day)
	body, err := p.fetch(ctx, url)
	if err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("olympics: fetch schedule: %w", err)
	}

	var resp apiScheduleResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return eventcore.SyncDelta{}, fmt.Errorf("olympics: decode schedule: %w", err)
	}

	mapped := mapSchedule(resp)
	delta := eventcore.SyncDelta{
		Competitions: mapped.competitions,
		Seasons:      mapped.seasons,
		Stages:       mapped.stages,
		Groups:       mapped.groups,
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

func (p Provider) fetch(ctx context.Context, url string) ([]byte, error) {
	cacheKey := "olympics_" + url
	if p.cache != nil {
		if data, ok, err := p.cache.Read(ctx, cacheKey); err == nil && ok {
			return data, nil
		}
	}

	resp, err := p.client.Do(ctx, httpdatasource.Request{
		Method: http.MethodGet,
		URL:    url,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("olympics: HTTP %d for %s", resp.StatusCode, url)
	}

	if p.cache != nil {
		_ = p.cache.Write(ctx, cacheKey, resp.Body, 0)
	}
	return resp.Body, nil
}
