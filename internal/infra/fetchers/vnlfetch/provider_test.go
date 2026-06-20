package vnlfetch

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jictyvoo/olhojogo/internal/domain/eventcore"
	"github.com/jictyvoo/olhojogo/internal/infra/httpdatasource"
)

const sampleScheduleJSON = `{
  "allTournaments": [
    {"no": 1661, "name": "Men's Volleyball Nations League 2026",
     "startDate": "2026-06-03T03:30:00Z", "endDate": "2026-08-02T20:00:00Z",
     "gender": "Men", "competitionSlug": "vnl-2026", "competitionFullName": "VNL 2026"},
    {"no": 1662, "name": "Women's Volleyball Nations League 2026",
     "startDate": "2026-06-04T03:30:00Z", "endDate": "2026-08-03T20:00:00Z",
     "gender": "Women", "competitionSlug": "vnl-2026", "competitionFullName": "VNL 2026"}
  ],
  "allTeams": [
    {"no": 8623, "code": "DOM", "country": "Dominican Republic", "name": "Dominican Republic"},
    {"no": 8634, "code": "USA", "country": "USA", "name": "USA"},
    {"no": 8590, "code": "BRA", "country": "Brazil", "name": "Brazil"},
    {"no": 8591, "code": "POL", "country": "Poland", "name": "Poland"}
  ],
  "matches": [
    {"matchNo": 26586, "matchDateUtc": "2026-06-17T12:00:00Z", "matchStatus": 2,
     "tournamentNo": 1662, "gender": "Women", "competitionSlug": "vnl-2026",
     "competitionFullName": "VNL 2026", "teamANo": 8623, "teamBNo": 8634,
     "winnerTeamNo": 8634, "teamAScore": 0, "teamBScore": 3,
     "roundNo": 297, "roundName": "Semana 2", "roundCode": "2",
     "pool": {"no": 5704, "name": "Grupo 5", "code": "5"},
     "city": "Pasig City", "countryCode": "PH", "isMatchTBD": false},
    {"matchNo": 26590, "matchDateUtc": "2099-06-20T16:30:00Z", "matchStatus": 0,
     "tournamentNo": 1661, "gender": "Men", "competitionSlug": "vnl-2026",
     "competitionFullName": "VNL 2026", "teamANo": 8590, "teamBNo": 8591,
     "winnerTeamNo": null, "teamAScore": -2147483648, "teamBScore": -2147483648,
     "roundNo": 297, "roundName": "Semana 2", "roundCode": "2",
     "pool": {"no": 5688, "name": "Grupo 4", "code": "4"},
     "city": "Ankara", "countryCode": "TR", "isMatchTBD": false},
    {"matchNo": 26999, "matchDateUtc": "2099-07-30T16:30:00Z", "matchStatus": 0,
     "tournamentNo": 1661, "gender": "Men", "competitionSlug": "vnl-2026",
     "competitionFullName": "VNL 2026", "teamANo": 0, "teamBNo": 0,
     "winnerTeamNo": null, "teamAScore": -2147483648, "teamBScore": -2147483648,
     "roundNo": 300, "roundName": "Finals", "roundCode": "3",
     "pool": {"no": 0, "name": "", "code": ""}, "isMatchTBD": true}
  ]
}`

func vnlTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleScheduleJSON))
	}))
}

func TestProvider_Identity(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "1661;1662")
	if p.Code() != eventcore.ProviderVNL {
		t.Errorf("Code = %q, want %q", p.Code(), eventcore.ProviderVNL)
	}
	if p.DisplayName() == "" {
		t.Error("DisplayName must not be empty")
	}
}

func TestProvider_SyncFixturesByDate_HappyPath(t *testing.T) {
	srv := vnlTestServer(t)
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "1661;1662")
	day := time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)

	delta, err := p.SyncFixturesByDate(t.Context(), day)
	if err != nil {
		t.Fatalf("SyncFixturesByDate: %v", err)
	}

	// The TBD final slot is dropped; two real fixtures remain.
	if len(delta.Fixtures) != 2 {
		t.Fatalf("expected 2 fixtures; got %d", len(delta.Fixtures))
	}
	// One competition shared by both genders, one season per tournament.
	if len(delta.Competitions) != 1 {
		t.Fatalf("expected 1 competition; got %d", len(delta.Competitions))
	}
	if len(delta.Seasons) != 2 {
		t.Fatalf("expected 2 seasons; got %d", len(delta.Seasons))
	}
	// roundNo 297 is shared across tournaments, so composite keys yield 2 stages.
	if len(delta.Stages) != 2 {
		t.Fatalf("expected 2 stages; got %d", len(delta.Stages))
	}
	if len(delta.Groups) != 2 {
		t.Fatalf("expected 2 groups; got %d", len(delta.Groups))
	}
	// Two distinct host cities -> two venues, carrying city + country code.
	if len(delta.Venues) != 2 {
		t.Fatalf("expected 2 venues; got %d", len(delta.Venues))
	}
	if len(delta.Participants) != 4 {
		t.Fatalf("expected 4 participants; got %d", len(delta.Participants))
	}
	// Only the finished match yields results (one per side).
	if len(delta.Results) != 2 {
		t.Fatalf("expected 2 results; got %d", len(delta.Results))
	}

	var finished eventcore.Fixture
	for _, f := range delta.Fixtures {
		if f.Status == eventcore.FixtureFinished {
			finished = f
		}
	}
	if finished.Name != "Dominican Republic vs USA" {
		t.Errorf("finished fixture name = %q", finished.Name)
	}
	if finished.GroupID == nil {
		t.Error("expected the finished fixture to carry a group")
	}
	if finished.VenueID == nil {
		t.Error("expected the finished fixture to carry a venue (host city)")
	}
}

func TestProvider_SyncFixturesByDate_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "1661;1662")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected error from 500 response")
	}
}

func TestProvider_SyncFixturesByDate_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{not json`))
	}))
	defer srv.Close()

	p := New(httpdatasource.New(), nil, srv.URL, "en", "1661;1662")
	if _, err := p.SyncFixturesByDate(t.Context(), time.Now()); err == nil {
		t.Fatal("expected JSON decode error")
	}
}

func TestProvider_New_AppliesDefaults(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "1661;1662")
	if p.baseURL == "" || p.lang == "" {
		t.Fatalf("expected defaults; got baseURL=%q lang=%q", p.baseURL, p.lang)
	}
}

func TestProvider_SyncFixtureResults_NotImplemented(t *testing.T) {
	p := New(httpdatasource.New(), nil, "", "", "1661;1662")
	if _, err := p.SyncFixtureResults(t.Context(), eventcore.Fixture{}); err == nil {
		t.Fatal("expected ErrNotImplemented")
	}
}
