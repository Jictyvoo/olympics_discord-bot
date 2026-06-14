//go:build curl

package httpdatasource

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/andelf/go-curl"
)

type curlClient struct{}

func New() Client { return &curlClient{} }

func (c *curlClient) Do(_ context.Context, req Request) (Response, error) {
	var buf bytes.Buffer
	easy := curl.EasyInit()
	defer easy.Cleanup()

	writeFunc := func(data []byte, _ any) bool {
		_, err := buf.Write(data)
		return err == nil
	}

	method := req.Method
	if method == "" {
		method = "GET"
	}

	errs := []error{
		easy.Setopt(curl.OPT_URL, req.URL),
		easy.Setopt(curl.OPT_HTTP_VERSION, curl.HTTP_VERSION_1_1),
		easy.Setopt(curl.OPT_CUSTOMREQUEST, method),
		easy.Setopt(curl.OPT_ENCODING, ""),
		easy.Setopt(curl.OPT_WRITEFUNCTION, writeFunc),
		easy.Setopt(curl.OPT_NOPROGRESS, true),
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, 0),
		easy.Setopt(curl.OPT_SSL_VERIFYHOST, 0),
	}
	if len(req.Headers) > 0 {
		hdrs := make([]string, 0, len(req.Headers))
		for k, v := range req.Headers {
			hdrs = append(hdrs, k+": "+v)
		}
		errs = append(errs, easy.Setopt(curl.OPT_HTTPHEADER, hdrs))
	}
	for _, err := range errs {
		if err != nil {
			return Response{}, fmt.Errorf("httpdatasource(curl): setup: %w", err)
		}
	}

	if err := easy.Perform(); err != nil {
		return Response{}, fmt.Errorf("httpdatasource(curl): perform: %w", err)
	}

	code, _ := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	statusCode, _ := code.(int)
	return Response{
		StatusCode: statusCode,
		Body:       buf.Bytes(),
	}, nil
}

func cacheKeyFromURL(method, url string) string {
	var b strings.Builder
	b.WriteString(method)
	b.WriteByte('-')
	for _, ch := range url {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '-' {
			b.WriteRune(ch)
		} else {
			b.WriteRune('_')
		}
	}
	b.WriteString(".log")
	return b.String()
}
