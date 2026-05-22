package schedule

import (
	"fmt"
	"time"
)

// NextN returns the next n execution times for the given Schedule after the
// provided start time. It advances minute-by-minute and checks Matches.
func NextN(s Schedule, after time.Time, n int) ([]time.Time, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be positive, got %d", n)
	}
	if n > 1000 {
		return nil, fmt.Errorf("n exceeds maximum of 1000")
	}

	results := make([]time.Time, 0, n)
	// Truncate to the next whole minute after 'after'
	current := after.Truncate(time.Minute).Add(time.Minute)

	// Safety limit: search up to 1 year ahead
	limit := after.Add(365 * 24 * time.Hour)

	for current.Before(limit) && len(results) < n {
		if Matches(s, current) {
			results = append(results, current)
		}
		current = current.Add(time.Minute)
	}

	if len(results) < n {
		return results, fmt.Errorf("only found %d occurrences within 1 year", len(results))
	}
	return results, nil
}

// Next returns the single next execution time for the Schedule after 'after'.
func Next(s Schedule, after time.Time) (time.Time, error) {
	times, err := NextN(s, after, 1)
	if err != nil {
		return time.Time{}, err
	}
	return times[0], nil
}
