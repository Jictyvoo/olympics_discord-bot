package utils

import "time"

func ParseTimestamp(timestamp string) (t time.Time, err error) {
	layouts := [...]string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700 -0700",
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err = time.Parse(layout, timestamp)
		if err == nil {
			return t, nil
		}
	}

	return time.Parse(time.DateTime, timestamp)
}

func EnsureTime(value *time.Time, duration time.Duration) time.Time {
	if value == nil {
		value = new(time.Time)
	}
	if !value.IsZero() {
		return *value
	}

	now := time.Now()
	*value = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		Add(duration)

	return *value
}
