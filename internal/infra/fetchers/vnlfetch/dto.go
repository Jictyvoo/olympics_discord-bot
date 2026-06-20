package vnlfetch

import "time"

type apiTournament struct {
	No                  int       `json:"no"`
	Name                string    `json:"name"`
	StartDate           time.Time `json:"startDate"`
	EndDate             time.Time `json:"endDate"`
	Gender              string    `json:"gender"`
	CompetitionSlug     string    `json:"competitionSlug"`
	CompetitionFullName string    `json:"competitionFullName"`
}

type apiTeam struct {
	No             int    `json:"no"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	TranslatedName string `json:"translatedName"`
	Country        string `json:"country"`
}

type apiPool struct {
	No   int    `json:"no"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type apiMatch struct {
	MatchNo             int       `json:"matchNo"`
	MatchDateUtc        time.Time `json:"matchDateUtc"`
	MatchStatus         int       `json:"matchStatus"`
	TournamentNo        int       `json:"tournamentNo"`
	Gender              string    `json:"gender"`
	CompetitionSlug     string    `json:"competitionSlug"`
	CompetitionFullName string    `json:"competitionFullName"`
	TeamANo             int       `json:"teamANo"`
	TeamBNo             int       `json:"teamBNo"`
	WinnerTeamNo        *int      `json:"winnerTeamNo"`
	TeamAScore          int       `json:"teamAScore"`
	TeamBScore          int       `json:"teamBScore"`
	RoundNo             int       `json:"roundNo"`
	RoundName           string    `json:"roundName"`
	RoundCode           string    `json:"roundCode"`
	Pool                apiPool   `json:"pool"`
	City                string    `json:"city"`
	CountryCode         string    `json:"countryCode"`
	IsMatchTBD          bool      `json:"isMatchTBD"`
}

type apiScheduleResponse struct {
	AllTournaments []apiTournament `json:"allTournaments"`
	AllTeams       []apiTeam       `json:"allTeams"`
	Matches        []apiMatch      `json:"matches"`
}
