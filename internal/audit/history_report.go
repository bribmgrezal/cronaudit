package audit

import (
	"fmt"
	"time"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// HistoryReport holds the history audit result for a single crontab file.
type HistoryReport struct {
	File    string
	From    time.Time
	Until   time.Time
	Entries []schedule.HistoryEntry
	Err     error
}

// AuditHistory parses a crontab file and collects all historical firing times
// for each entry within [from, until).
func AuditHistory(path string, from, until time.Time) HistoryReport {
	report := HistoryReport{File: path, From: from, Until: until}

	entries, err := parser.Parse(path)
	if err != nil {
		report.Err = err
		return report
	}

	for _, entry := range entries {
		hist, err := schedule.BuildHistory(entry.Schedule, entry.Command, from, until)
		if err != nil {
			// skip unparseable schedules but continue
			continue
		}
		report.Entries = append(report.Entries, hist...)
	}
	return report
}

// FormatHistoryReport returns a human-readable string for a HistoryReport.
func FormatHistoryReport(r HistoryReport) string {
	if r.Err != nil {
		return fmt.Sprintf("error reading %s: %v\n", r.File, r.Err)
	}
	header := fmt.Sprintf("History for %s (%s – %s):\n",
		r.File,
		r.From.UTC().Format("2006-01-02 15:04"),
		r.Until.UTC().Format("2006-01-02 15:04"),
	)
	return header + schedule.FormatHistory(r.Entries)
}
