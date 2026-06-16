package fifafetch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

// fetchSeason is best-effort: on any failure it falls back to the requested day
// so the not-null timestamp columns still get a valid span.
func (p Provider) fetchSeason(ctx context.Context, day time.Time) seasonMeta {
	fallback := seasonMeta{startsOn: day.UTC(), endsOn: day.UTC()}
	if p.seasonID == "" {
		return fallback
	}

	body, err := p.fetch(ctx, seasonURL(p.baseURL, p.seasonID, p.lang))
	if err != nil {
		return fallback
	}
	var s apiSeason
	if err = json.Unmarshal(body, &s); err != nil {
		return fallback
	}

	meta := seasonMeta{
		name:     localized(s.Name, p.lang),
		startsOn: s.StartDate.UTC(),
		endsOn:   s.EndDate.UTC(),
	}
	if meta.startsOn.IsZero() {
		meta.startsOn = day.UTC()
	}
	if meta.endsOn.IsZero() {
		meta.endsOn = meta.startsOn
	}
	return meta
}

// fetchStandings joins per-stage errors without dropping the standings that succeeded.
func (p Provider) fetchStandings(
	ctx context.Context,
	stageKeys []string,
) ([]eventcore.Standing, error) {
	var (
		out  []eventcore.Standing
		errs []error
	)
	for _, stageKey := range stageKeys {
		url := standingURL(p.baseURL, p.competitionID, p.seasonID, stageKey, p.lang)
		body, err := p.fetch(ctx, url)
		if err != nil {
			errs = append(errs, fmt.Errorf("fifa: fetch standing %s: %w", stageKey, err))
			continue
		}
		var resp apiStandingResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			errs = append(errs, fmt.Errorf("fifa: decode standing %s: %w", stageKey, err))
			continue
		}
		out = append(out, mapStandings(resp)...)
	}
	return out, errors.Join(errs...)
}

func (p Provider) fetch(ctx context.Context, url string) ([]byte, error) {
	cacheKey := "fifa_" + url
	if p.cache != nil {
		if data, ok, err := p.cache.Read(ctx, cacheKey); err == nil && ok {
			return data, nil
		}
	}

	resp, err := p.client.Do(ctx, httpdatasource.Request{
		Method:  http.MethodGet,
		URL:     url,
		Headers: map[string]string{"Accept": "application/json"},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fifa: HTTP %d for %s", resp.StatusCode, url)
	}

	if p.cache != nil {
		_ = p.cache.Write(ctx, cacheKey, resp.Body, 0)
	}
	return resp.Body, nil
}
