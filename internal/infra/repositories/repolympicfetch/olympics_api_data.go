package repolympicfetch

import "time"

type OlympicsAPIResponseCompetitors struct {
	Code    string `json:"code"`
	Noc     string `json:"noc,omitempty"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Results struct {
		Position  string `json:"position"`
		Mark      string `json:"mark"`
		MedalType string `json:"medalType"`
		Irm       string `json:"irm"`
	} `json:"results,omitempty"`
}

type (
	OlympicsAPIResponseUnit struct {
		DisciplineName      string                           `json:"disciplineName"`
		EventUnitName       string                           `json:"eventUnitName"`
		Id                  string                           `json:"id"`
		DisciplineCode      string                           `json:"disciplineCode"`
		GenderCode          string                           `json:"genderCode"`
		EventCode           string                           `json:"eventCode"`
		PhaseCode           string                           `json:"phaseCode"`
		EventId             string                           `json:"eventId"`
		EventName           string                           `json:"eventName"`
		PhaseId             string                           `json:"phaseId"`
		PhaseName           string                           `json:"phaseName"`
		DisciplineId        string                           `json:"disciplineId"`
		EventOrder          int                              `json:"eventOrder"`
		PhaseType           string                           `json:"phaseType"`
		EventUnitType       string                           `json:"eventUnitType"`
		OlympicDay          string                           `json:"olympicDay"`
		StartDate           time.Time                        `json:"startDate"`
		EndDate             time.Time                        `json:"endDate"`
		HideStartDate       bool                             `json:"hideStartDate"`
		HideEndDate         bool                             `json:"hideEndDate"`
		StartText           string                           `json:"startText"`
		Order               int                              `json:"order"`
		Venue               string                           `json:"venue"`
		VenueDescription    string                           `json:"venueDescription"`
		Location            string                           `json:"location"`
		LocationDescription string                           `json:"locationDescription"`
		Status              string                           `json:"status"`
		StatusDescription   string                           `json:"statusDescription"`
		MedalFlag           int                              `json:"medalFlag"`
		LiveFlag            bool                             `json:"liveFlag"`
		ScheduleItemType    string                           `json:"scheduleItemType"`
		UnitNum             string                           `json:"unitNum"`
		SessionCode         string                           `json:"sessionCode"`
		Competitors         []OlympicsAPIResponseCompetitors `json:"competitors"`
		ExtraData           struct {
			DetailUrl string `json:"detailUrl"`
		} `json:"extraData"`
		GroupId string `json:"groupId,omitempty"`
	}
	OlympicsAPIResponse struct {
		Units  []OlympicsAPIResponseUnit `json:"units"`
		Groups []any                     `json:"groups"`
	}
)
