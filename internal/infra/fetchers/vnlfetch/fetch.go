package vnlfetch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

// browserUserAgent keeps the site-XHR request looking like the page.
const browserUserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:152.0) Gecko/20100101 Firefox/152.0"

func (p Provider) fetch(ctx context.Context, url string) ([]byte, error) {
	cacheKey := "vnl_" + url
	if p.cache != nil {
		if data, ok, err := p.cache.Read(ctx, cacheKey); err == nil && ok {
			return data, nil
		}
	}

	resp, err := p.client.Do(ctx, httpdatasource.Request{
		Method: http.MethodGet,
		URL:    url,
		Headers: map[string]string{
			"Accept":           "application/json",
			"User-Agent":       browserUserAgent,
			"X-Requested-With": "XMLHttpRequest",
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vnl: HTTP %d for %s", resp.StatusCode, url)
	}

	if p.cache != nil {
		_ = p.cache.Write(ctx, cacheKey, resp.Body, 0)
	}
	return resp.Body, nil
}
