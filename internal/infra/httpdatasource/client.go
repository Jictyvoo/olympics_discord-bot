package httpdatasource

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Request struct {
	Method   string
	URL      string
	Headers  map[string]string
	Body     []byte
	Timeout  time.Duration
	CacheKey string
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Client interface {
	Do(ctx context.Context, req Request) (Response, error)
}

const (
	defaultMaxIdleConns        = 100
	defaultIdleConnTimeoutSecs = 90
	defaultMaxIdleConnsPerHost = 10
	defaultRequestTimeoutSecs  = 30
)

type httpClient struct {
	client *http.Client
}

//nolint:ireturn // factory returning consumer interface by design
func New() Client {
	transport := &http.Transport{
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        defaultMaxIdleConns,
		IdleConnTimeout:     defaultIdleConnTimeoutSecs * time.Second,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &httpClient{
		client: &http.Client{Transport: transport},
	}
}

func (c *httpClient) Do(ctx context.Context, req Request) (Response, error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = defaultRequestTimeoutSecs * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	method := req.Method
	if method == "" {
		method = http.MethodGet
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, nil)
	if err != nil {
		return Response{}, fmt.Errorf("httpdatasource: build request: %w", err)
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("httpdatasource: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("httpdatasource: read body: %w", err)
	}
	return Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
