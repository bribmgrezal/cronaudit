package schedule

import (
	"time"
)

// MissedWindow represents a time range during which a schedule should have fired
// but did not (e.g., due to downtime).
type MissedWindow struct {
	Start time.Time
	End   time.Time
}

// FindMissed returns all times within [window.Start, window.End) at which the
// given schedule would have fired. Granularity is per-minute.
func FindMissed(s *Schedule, window MissedWindow) []time.Time {
	var missed []time.Time

	// Truncate to minute boundary
	current := window.Start.Truncate(time.Minute)
	end := window.End.Truncate(time.Minute)

	for !current.After(end) {
		if Matches(s, current) {
			missed = append(missed, current)
		}
		current = current.Add(time.Minute)
	}

	return missed
}

// Matches reports whether the schedule would fire at the given time.
func Matches(s *Schedule, t time.Time) bool {
	return contains(s.Minutes, t.Minute()) &&
		contains(s.Hours, t.Hour()) &&
		contains(s.Days, t.Day()) &&
		contains(s.Months, int(t.Month())) &&
		contains(s.Weekdays, int(t.Weekday()))
}

func contains(set []int, val int) bool {
	for _, v := range set {
		if v == val {
			return true
		}
	}
	return false
}
