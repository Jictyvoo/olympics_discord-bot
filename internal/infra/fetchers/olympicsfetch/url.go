package olympicsfetch

import (
	"fmt"
	"time"
)

func scheduleByDayURL(baseURL, lang string, day time.Time) string {
	return fmt.Sprintf(
		"%s/summer/schedules/api/%s/schedule/day/%s",
		baseURL,
		lang,
		day.Format(time.DateOnly),
	)
}
