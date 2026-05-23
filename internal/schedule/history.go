package schedule

import (
	"fmt"
	"time"
)

// HistoryEntry represents a single scheduled execution in the past.
type HistoryEntry struct {
	Time       time.Time
	Expression string
	Command    string
}

// BuildHistory returns all times a cron expression would have fired
// within the half-open interval [from, until).
func BuildHistory(expr, command string, from, until time.Time) ([]HistoryEntry, error) {
	sched, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", expr, err)
	}

	var entries []HistoryEntry
	// Align start to the next whole minute at or after from.
	current := from.Truncate(time.Minute)
	if current.Before(from) {
		current = current.Add(time.Minute)
	}

	for !current.After(until) && !current.Equal(until) {
		if Matches(sched, current) {
			entries = append(entries, HistoryEntry{
				Time:       current,
				Expression: expr,
				Command:    command,
			})
		}
		current = current.Add(time.Minute)
	}
	return entries, nil
}

// FormatHistory returns a human-readable summary of history entries.
func FormatHistory(entries []HistoryEntry) string {
	if len(entries) == 0 {
		return "no executions found in window\n"
	}
	out := fmt.Sprintf("%d execution(s) found:\n", len(entries))
	for _, e := range entries {
		out += fmt.Sprintf("  %s  %s  %s\n",
			e.Time.UTC().Format("2006-01-02 15:04"),
			e.Expression,
			e.Command,
		)
	}
	return out
}
