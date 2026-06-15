package olympicsfetch

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

const sampleScheduleJSON = `{
  "units": [
    {
      "id": "u-100m-final",
      "disciplineName": "Athletics",
      "eventUnitName": "100m Final",
      "eventName": "100m",
      "phaseName": "Final",
      "startDate": "2024-08-04T20:55:00Z",
      "endDate":   "2024-08-04T21:05:00Z",
      "status": "Finished",
      "competitors": [
        {"code": "USA-001", "noc": "USA", "name": "A. Sprinter", "order": 1,
         "results": {"position": "1", "mark": "9.79", "medalType": "ME_GOLD"}},
        {"code": "JAM-002", "noc": "JAM", "name": "B. Quickstep", "order": 2,
         "results": {"position": "2", "mark": "9.81", "medalType": "ME_SILVER"}}
      ]
    }
  ],
  "groups": []
}`

func TestProvider_SyncFixturesByDate_HappyPath(t *testing.T) {
	var gotPath, gotAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAccept = r.Header.Get("Accept")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleScheduleJSON))
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "ENG")
	day := time.Date(2024, 8, 4, 0, 0, 0, 0, time.UTC)

	delta, err := p.SyncFixturesByDate(t.Context(), day)
	if err != nil {
		t.Fatalf("SyncFixturesByDate: %v", err)
	}
	if !strings.Contains(gotPath, "/summer/schedules/api/ENG/schedule/day/2024-08-04") {
		t.Errorf("URL path = %q, expected to contain language + date segments", gotPath)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q, want application/json", gotAccept)
	}
	if len(delta.Fixtures) != 1 {
		t.Fatalf("expected 1 fixture; got %d", len(delta.Fixtures))
	}
	if delta.Fixtures[0].Status != eventcore.FixtureFinished {
		t.Errorf(
			"fixture status = %q, want %q",
			delta.Fixtures[0].Status,
			eventcore.FixtureFinished,
		)
	}
	if len(delta.Participants) != 2 {
		t.Fatalf("expected 2 participants; got %d", len(delta.Participants))
	}
}

func TestProvider_SyncFixturesByDate_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "ENG")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected error from 500 response")
	}
}

func TestProvider_SyncFixturesByDate_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "ENG")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected JSON decode error")
	}
}

func TestProvider_New_AppliesDefaults(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "")
	if p.baseURL == "" || p.lang == "" {
		t.Fatalf("expected defaults to be applied; got baseURL=%q lang=%q", p.baseURL, p.lang)
	}
}
