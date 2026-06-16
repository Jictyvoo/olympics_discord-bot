package fifafetch

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

const sampleMatchesJSON = `{
  "Results": [
    {
      "IdMatch": "400021443",
      "IdCompetition": "17",
      "IdSeason": "285023",
      "IdStage": "289273",
      "IdGroup": "289275",
      "StageName": [{"Locale": "en-GB", "Description": "First Stage"}],
      "GroupName": [{"Locale": "en-GB", "Description": "Group A"}],
      "CompetitionName": [{"Locale": "en-GB", "Description": "FIFA World Cup"}],
      "SeasonName": [{"Locale": "en-GB", "Description": "FIFA World Cup 2026"}],
      "Date": "2026-06-11T19:00:00Z",
      "Home": {"IdTeam": "43911", "IdCountry": "MEX", "Abbreviation": "MEX", "Gender": 1,
               "TeamName": [{"Locale": "en-GB", "Description": "Mexico"}]},
      "Away": {"IdTeam": "43883", "IdCountry": "RSA", "Abbreviation": "RSA", "Gender": 1,
               "TeamName": [{"Locale": "en-GB", "Description": "South Africa"}]},
      "HomeTeamScore": 2,
      "AwayTeamScore": 0,
      "Winner": "43911",
      "Stadium": {"IdStadium": "400222084", "IdCountry": "MEX",
                  "Name": [{"Locale": "en-GB", "Description": "Mexico City Stadium"}],
                  "CityName": [{"Locale": "en-GB", "Description": "Mexico City"}]},
      "MatchStatus": 0
    }
  ]
}`

const sampleStandingJSON = `{
  "Results": [
    {"IdStage": "289273", "IdGroup": "289275", "Position": 1, "Points": 3,
     "Won": 1, "Lost": 0, "Drawn": 0, "Played": 1, "For": 2, "Against": 0, "GoalsDiference": 2,
     "Team": {"IdTeam": "43911", "IdCountry": "MEX", "Abbreviation": "MEX",
              "Name": [{"Locale": "en-GB", "Description": "Mexico"}]}}
  ]
}`

func fifaTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/standing") {
			_, _ = w.Write([]byte(sampleStandingJSON))
			return
		}
		_, _ = w.Write([]byte(sampleMatchesJSON))
	}))
}

func TestProvider_Identity(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "17", "285023")
	if p.Code() != eventcore.ProviderFIFA {
		t.Errorf("Code = %q, want %q", p.Code(), eventcore.ProviderFIFA)
	}
	if p.DisplayName() == "" {
		t.Error("DisplayName must not be empty")
	}
}

func TestProvider_SyncFixturesByDate_HappyPath(t *testing.T) {
	srv := fifaTestServer(t)
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "17", "285023")
	day := time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC)

	delta, err := p.SyncFixturesByDate(t.Context(), day)
	if err != nil {
		t.Fatalf("SyncFixturesByDate: %v", err)
	}
	if len(delta.Fixtures) != 1 {
		t.Fatalf("expected 1 fixture; got %d", len(delta.Fixtures))
	}
	fx := delta.Fixtures[0]
	if fx.Status != eventcore.FixtureFinished {
		t.Errorf("fixture status = %q, want finished", fx.Status)
	}
	if fx.Name != "Mexico vs South Africa" {
		t.Errorf("fixture name = %q", fx.Name)
	}
	if fx.VenueID == nil {
		t.Error("expected venue to be set")
	}
	if len(delta.Participants) != 2 {
		t.Fatalf("expected 2 participants; got %d", len(delta.Participants))
	}
	if len(delta.Results) != 2 {
		t.Fatalf("expected 2 results; got %d", len(delta.Results))
	}
	if len(delta.Standings) != 1 {
		t.Fatalf("expected 1 standing; got %d", len(delta.Standings))
	}
	if delta.Standings[0].Rank != 1 || delta.Standings[0].Points != 3 {
		t.Errorf("unexpected standing: %+v", delta.Standings[0])
	}
}

func TestProvider_SyncFixturesByDate_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "17", "285023")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected error from 500 response")
	}
}

func TestProvider_SyncFixturesByDate_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "17", "285023")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected JSON decode error")
	}
}

func TestProvider_New_AppliesDefaults(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "17", "285023")
	if p.baseURL == "" || p.lang == "" {
		t.Fatalf("expected defaults; got baseURL=%q lang=%q", p.baseURL, p.lang)
	}
}

func TestProvider_SyncFixtureResults_NotImplemented(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "17", "285023")
	if _, err := p.SyncFixtureResults(t.Context(), eventcore.Fixture{}); err == nil {
		t.Fatal("expected ErrNotImplemented")
	}
}
