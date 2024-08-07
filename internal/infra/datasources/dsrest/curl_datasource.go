package dsrest

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"unicode"

	"github.com/andelf/go-curl"

	"github.com/jictyvoo/olympics_data_fetcher/internal/infra/datasources"
)

type responseInterceptor struct {
	bytes.Buffer
	headersBuffer bytes.Buffer
}

func (ri *responseInterceptor) Interceptor(buf []byte, _ any) bool {
	_, err := ri.Write(buf)
	return err == nil
}

func (ri *responseInterceptor) HeaderFunction(buf []byte, _ any) bool {
	_, err := ri.headersBuffer.Write(buf)
	return err == nil
}

type CurlDatasource struct {
	disableSSL bool
	cacher     datasources.CacheableDataSource
}

func NewCurlDatasource(disableSSL bool) RESTDataSource {
	return CurlDatasource{disableSSL: disableSSL}
}

func NewCurlCacheableDatasource(cache datasources.CacheableDataSource) RESTDataSource {
	return CurlDatasource{disableSSL: true, cacher: cache}
}

func (ds CurlDatasource) cacheKey(url string, method HTTPMethod) string {
	var builder strings.Builder
	builder.WriteString(string(method))
	builder.WriteByte('-')
	for _, character := range url {
		if unicode.IsLetter(character) || unicode.IsNumber(character) || character == '-' {
			builder.WriteRune(character)
		} else {
			builder.WriteRune('_')
		}
	}

	builder.WriteString(".log")
	return builder.String()
}

func (ds CurlDatasource) retrieveCache(url string, method HTTPMethod) []byte {
	if ds.cacher == nil {
		return nil
	}

	cacheKey := ds.cacheKey(url, method)
	found, err := ds.cacher.Read(cacheKey)
	if err != nil {
		return nil
	}
	return found
}

func (ds CurlDatasource) saveOnCache(url string, method HTTPMethod, body []byte) error {
	if ds.cacher == nil {
		return nil
	}

	cacheKey := ds.cacheKey(url, method)
	return ds.cacher.Write(cacheKey, body)
}

func (ds CurlDatasource) makeRequest(url string, method HTTPMethod) (HTTPResponse, error) {
	var responseBody responseInterceptor
	if cachedBody := ds.retrieveCache(url, method); cachedBody != nil {
		responseBody.Write(cachedBody)
		return HTTPResponse{
			StatusCode: 0,
			Body:       &responseBody,
			Headers:    map[string][]string{},
		}, nil
	}

	easy := curl.EasyInit()
	defer easy.Cleanup()

	var writeFunc any
	if method != MethodHead {
		writeFunc = responseBody.Interceptor
	}

	setupErrList := []error{
		easy.Setopt(curl.OPT_URL, url),
		easy.Setopt(curl.OPT_HTTP_VERSION, curl.HTTP_VERSION_1_1),
		easy.Setopt(curl.OPT_CUSTOMREQUEST, string(method)),
		easy.Setopt(curl.OPT_ENCODING, ""),
		easy.Setopt(curl.OPT_HEADERFUNCTION, responseBody.HeaderFunction),
		easy.Setopt(curl.OPT_WRITEFUNCTION, writeFunc),
		easy.Setopt(curl.OPT_NOPROGRESS, true),
		easy.Setopt(curl.OPT_VERBOSE, true),
	}

	if method == MethodHead {
		setupErrList = append(
			setupErrList, easy.Setopt(curl.OPT_NOBODY, true),
		)
	}

	if ds.disableSSL {
		setupErrList = append(
			setupErrList,
			easy.Setopt(curl.OPT_SSL_VERIFYPEER, 0), // Disable SSL verification
			easy.Setopt(curl.OPT_SSL_VERIFYHOST, 0), // Disable SSL host verification
		)
	}

	if err := errors.Join(setupErrList...); err != nil {
		return HTTPResponse{}, err
	}

	err := easy.Perform()
	if err != nil {
		slog.Error(
			"Failed to perform http call",
			slog.String("error", err.Error()),
			slog.String("url", url),
		)
	}

	respCode, _ := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	respObj := HTTPResponse{
		Body:    &responseBody,
		Headers: map[string][]string{},
	}
	respObj.StatusCode, _ = respCode.(int)
	// Save before returning
	if err == nil && responseBody.Len() > 0 {
		err = ds.saveOnCache(url, method, responseBody.Bytes())
	}
	return respObj, err
}

func (ds CurlDatasource) Get(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodGet)
}

func (ds CurlDatasource) Head(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodHead)
}

func (ds CurlDatasource) Delete(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodDelete)
}

func (ds CurlDatasource) Post(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodPost)
}

func (ds CurlDatasource) Put(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodPut)
}

func (ds CurlDatasource) Patch(url string) (HTTPResponse, error) {
	return ds.makeRequest(url, MethodPatch)
}
