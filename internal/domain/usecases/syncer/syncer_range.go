package syncer

import (
	"errors"
	"time"
)

const hoursPerDay = 24

// SyncRange fetches and persists every day in [from, to] inclusive, advancing
// one calendar day at a time. A failing day is recorded and the range
// continues; all day errors are joined into the result.
func (s *Syncer) SyncRange(from, to time.Time) error {
	from = from.UTC().Truncate(hoursPerDay * time.Hour)
	to = to.UTC().Truncate(hoursPerDay * time.Hour)

	var errs []error
	for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
		if err := s.SyncDay(day); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
