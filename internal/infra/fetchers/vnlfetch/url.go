package vnlfetch

import (
	"fmt"
	"time"
)

// scheduleByDayURL queries one UTC day; the feed filters inclusively, so from
// and to are both the day. tournaments is the ";"-joined list (e.g. "1661;1662").
func scheduleByDayURL(baseURL, tournaments string, day time.Time) string {
	d := day.UTC().Format(time.DateOnly)
	return fmt.Sprintf("%s/api/v1/volley-tournament/%s/%s/%s", baseURL, d, d, tournaments)
}
