package dsrest

import "io"

type HTTPMethod string

// Common HTTP methods.
//
// Unless otherwise noted, these are defined in RFC 7231 section 4.3.
const (
	MethodGet     HTTPMethod = "GET"
	MethodHead    HTTPMethod = "HEAD"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodPatch   HTTPMethod = "PATCH" // RFC 5789
	MethodDelete  HTTPMethod = "DELETE"
	MethodConnect HTTPMethod = "CONNECT"
	MethodOptions HTTPMethod = "OPTIONS"
	MethodTrace   HTTPMethod = "TRACE"
)

type HTTPResponse struct {
	StatusCode int
	Body       io.Reader
	Headers    map[string][]string
}

type RESTDataSource interface {
	Get(url string) (HTTPResponse, error)
	Head(url string) (HTTPResponse, error)
	Delete(url string) (HTTPResponse, error)
	Post(url string) (HTTPResponse, error)
	Put(url string) (HTTPResponse, error)
	Patch(url string) (HTTPResponse, error)
}
