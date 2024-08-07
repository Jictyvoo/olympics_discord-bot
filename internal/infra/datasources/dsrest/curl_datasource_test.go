package dsrest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCurlDatasource(t *testing.T) {
	tests := []struct {
		name           string
		method         HTTPMethod
		handler        http.HandlerFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Get",
			method: MethodGet, expectedStatus: http.StatusOK,
			expectedBody: "Get response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Get response"))
			},
		},
		{
			name: "Head", method: MethodHead, expectedStatus: http.StatusOK,
			expectedBody: "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "Delete", method: MethodDelete, expectedStatus: http.StatusNoContent,
			expectedBody: "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
		{
			name: "Post", method: MethodPost, expectedStatus: http.StatusCreated,
			expectedBody: "Post response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("Post response"))
			},
		},
		{
			name: "Put", method: MethodPut, expectedStatus: http.StatusOK,
			expectedBody: "Put response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Put response"))
			},
		},
		{
			name: "Patch", method: MethodPatch, expectedStatus: http.StatusOK,
			expectedBody: "Patch response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Patch response"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handler := tt.handler
				if handler == nil {
					handler = func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}

				server := httptest.NewServer(handler)
				defer server.Close()

				ds := NewCurlDatasource(true)

				var (
					resp HTTPResponse
					err  error
				)
				switch tt.method {
				case MethodGet:
					resp, err = ds.Get(server.URL)
				case MethodHead:
					resp, err = ds.Head(server.URL)
				case MethodDelete:
					resp, err = ds.Delete(server.URL)
				case MethodPost:
					resp, err = ds.Post(server.URL)
				case MethodPut:
					resp, err = ds.Put(server.URL)
				case MethodPatch:
					resp, err = ds.Patch(server.URL)
				}

				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.StatusCode != tt.expectedStatus {
					t.Errorf("expected status %v, got %v", tt.expectedStatus, resp.StatusCode)
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if string(bodyBytes) != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, resp.Body)
				}
			},
		)
	}
}
