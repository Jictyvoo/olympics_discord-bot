package fifafetch

import "time"

// localizedText is the upstream shape for every translatable string.
type localizedText []struct {
	Locale      string
	Description string
}

type apiTeam struct {
	IdTeam       string
	TeamName     localizedText
	Abbreviation string
	IdCountry    string
	Gender       int
	Score        *int
}

type apiStadium struct {
	IdStadium string
	Name      localizedText
	CityName  localizedText
	IdCountry string
}

type apiMatch struct {
	IdMatch         string
	IdCompetition   string
	IdSeason        string
	IdStage         string
	IdGroup         string
	StageName       localizedText
	GroupName       localizedText
	CompetitionName localizedText
	SeasonName      localizedText
	Date            time.Time
	Home            apiTeam
	Away            apiTeam
	HomeTeamScore   *int
	AwayTeamScore   *int
	Winner          string
	Stadium         apiStadium
	MatchStatus     int
}

type apiMatchesResponse struct {
	Results []apiMatch
}

type apiSeason struct {
	IdSeason  string
	Name      localizedText
	StartDate time.Time
	EndDate   time.Time
}

type apiStandingTeam struct {
	IdTeam       string
	Name         localizedText
	IdCountry    string
	Abbreviation string
}

type apiStanding struct {
	IdStage        string
	IdGroup        string
	Position       int
	Points         int
	Won            int
	Lost           int
	Drawn          int
	Played         int
	For            int
	Against        int
	GoalsDiference int
	Team           apiStandingTeam
}

type apiStandingResponse struct {
	Results []apiStanding
}
