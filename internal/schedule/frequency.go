package schedule

import "time"

// Frequency describes how often a schedule fires.
type Frequency struct {
	// EstimatedPerDay is the approximate number of times the schedule fires in a 24-hour period.
	EstimatedPerDay float64
	// Label is a human-readable description of the frequency.
	Label string
}

// ClassifyFrequency returns a Frequency for the given Schedule based on how
// many times it would fire in a standard 24-hour window.
func ClassifyFrequency(s Schedule) Frequency {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	count := CountInWindow(s, start, end)

	var label string
	switch {
	case count == 0:
		label = "never (in a typical day)"
	case count == 1:
		label = "daily"
	case count == 2:
		label = "twice daily"
	case count <= 4:
		label = "a few times daily"
	case count == 24:
		label = "hourly"
	case count > 24 && count <= 60:
		label = "multiple times per hour"
	case count == 1440:
		label = "every minute"
	case count > 1440:
		label = "sub-minute (unusual)"
	default:
		label = "custom"
	}

	return Frequency{
		EstimatedPerDay: float64(count),
		Label:           label,
	}
}

// IsHighFrequency returns true if the schedule fires more than once per hour on average.
func IsHighFrequency(s Schedule) bool {
	f := ClassifyFrequency(s)
	return f.EstimatedPerDay > 24
}
