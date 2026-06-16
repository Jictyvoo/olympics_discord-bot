package fifafetch

import (
	"fmt"
	"net/url"
	"time"
)

const matchPageSize = "200"

// matchesByDayURL scopes the calendar/matches query to a single UTC day.
func matchesByDayURL(baseURL, competitionID, seasonID, lang string, day time.Time) string {
	from := day.UTC().Truncate(hoursPerDay * time.Hour)
	to := from.Add(hoursPerDay * time.Hour)

	q := url.Values{}
	q.Set("idCompetition", competitionID)
	q.Set("idSeason", seasonID)
	q.Set("from", from.Format(time.RFC3339))
	q.Set("to", to.Format(time.RFC3339))
	q.Set("count", matchPageSize)
	q.Set("language", lang)

	return baseURL + "/calendar/matches?" + q.Encode()
}

func seasonURL(baseURL, seasonID, lang string) string {
	q := url.Values{}
	q.Set("language", lang)
	return fmt.Sprintf("%s/seasons/%s?%s", baseURL, seasonID, q.Encode())
}

func standingURL(baseURL, competitionID, seasonID, stageKey, lang string) string {
	q := url.Values{}
	q.Set("language", lang)
	q.Set("count", matchPageSize)

	return fmt.Sprintf(
		"%s/calendar/%s/%s/%s/standing?%s",
		baseURL, competitionID, seasonID, stageKey, q.Encode(),
	)
}

const hoursPerDay = 24
