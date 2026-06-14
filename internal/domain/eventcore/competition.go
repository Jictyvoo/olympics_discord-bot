package eventcore

import "time"

type Competition struct {
	ID         CanonicalID
	Ext        ExternalID
	Code       string
	Name       string
	Discipline string
}

type Season struct {
	ID            CanonicalID
	Ext           ExternalID
	CompetitionID CanonicalID
	Name          string
	StartsOn      time.Time
	EndsOn        time.Time
}

type Stage struct {
	ID       CanonicalID
	Ext      ExternalID
	SeasonID CanonicalID
	Name     string
	Ord      int
}

type Group struct {
	ID      CanonicalID
	Ext     ExternalID
	StageID CanonicalID
	Name    string
}
