package schedule

import "time"

// Window represents a time range with a start and end.
type Window struct {
	Start time.Time
	End   time.Time
}

// NextInWindow returns all scheduled times within the given window for a
// parsed schedule. It uses Next repeatedly to walk forward through the window.
func NextInWindow(s *Schedule, w Window) []time.Time {
	var results []time.Time

	if w.End.Before(w.Start) {
		return results
	}

	current := w.Start
	for {
		next := Next(s, current)
		if next.IsZero() || next.After(w.End) {
			break
		}
		results = append(results, next)
		// Advance by one minute to avoid returning the same time.
		current = next.Add(time.Minute)
		if current.After(w.End) {
			break
		}
	}

	return results
}

// CountInWindow returns the number of times a schedule fires within the window.
func CountInWindow(s *Schedule, w Window) int {
	return len(NextInWindow(s, w))
}
