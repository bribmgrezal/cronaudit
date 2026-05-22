package schedule

import "sort"

// Conflict describes two schedule entries that overlap in time.
type Conflict struct {
	HostA    string
	CommandA string
	HostB    string
	CommandB string
}

// Entry pairs a host and command with its parsed schedule.
type Entry struct {
	Host     string
	Command  string
	Schedule *Schedule
}

// FindConflicts returns pairs of entries whose schedules share at least one
// common trigger time (same minute, hour, day-of-month, month, day-of-week).
func FindConflicts(entries []Entry) []Conflict {
	var conflicts []Conflict
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if schedulesOverlap(entries[i].Schedule, entries[j].Schedule) {
				conflicts = append(conflicts, Conflict{
					HostA:    entries[i].Host,
					CommandA: entries[i].Command,
					HostB:    entries[j].Host,
					CommandB: entries[j].Command,
				})
			}
		}
	}
	return conflicts
}

// schedulesOverlap returns true if two schedules share any common firing time.
func schedulesOverlap(a, b *Schedule) bool {
	return hasCommon(a.Minute.Values, b.Minute.Values) &&
		hasCommon(a.Hour.Values, b.Hour.Values) &&
		hasCommon(a.DayOfMonth.Values, b.DayOfMonth.Values) &&
		hasCommon(a.Month.Values, b.Month.Values) &&
		hasCommon(a.DayOfWeek.Values, b.DayOfWeek.Values)
}

// hasCommon returns true if two sorted integer slices share at least one value.
func hasCommon(a, b []int) bool {
	if !sort.IntsAreSorted(a) {
		sort.Ints(a)
	}
	if !sort.IntsAreSorted(b) {
		sort.Ints(b)
	}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] == b[j]:
			return true
		case a[i] < b[j]:
			i++
		default:
			j++
		}
	}
	return false
}
