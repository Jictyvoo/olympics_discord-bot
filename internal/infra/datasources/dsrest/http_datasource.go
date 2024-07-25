package dsrest

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HTTPDatasource struct {
	client *http.Client
}

func NewHTTPDatasource() *HTTPDatasource {
	client := &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: false, // Force HTTP/1.1
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Disable SSL verification
			},
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 10,
		},
		Timeout: 10 * time.Second, // Set timeout
	}
	return &HTTPDatasource{client: client}
}
