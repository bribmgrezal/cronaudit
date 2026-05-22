package audit

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// ConflictReport holds the results of auditing crontab entries.
type ConflictReport struct {
	Host      string
	Conflicts []schedule.Conflict
	Errors    []string
}

// Audit parses a crontab file and finds scheduling conflicts.
func Audit(host string, r io.Reader) (*ConflictReport, error) {
	entries, err := parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parsing crontab for host %s: %w", host, err)
	}

	report := &ConflictReport{Host: host}

	var schedules []schedule.Schedule
	for _, entry := range entries {
		sched, err := schedule.Parse(entry.Schedule)
		if err != nil {
			report.Errors = append(report.Errors,
				fmt.Sprintf("line %d: invalid schedule %q: %v", entry.Line, entry.Schedule, err))
			continue
		}
		sched.Command = entry.Command
		sched.Line = entry.Line
		schedules = append(schedules, sched)
	}

	report.Conflicts = schedule.FindConflicts(schedules)
	return report, nil
}

// Format writes a human-readable summary of the report to w.
func (r *ConflictReport) Format(w io.Writer) {
	fmt.Fprintf(w, "=== Audit Report: %s ===\n", r.Host)

	if len(r.Errors) > 0 {
		fmt.Fprintln(w, "\nParse Errors:")
		for _, e := range r.Errors {
			fmt.Fprintf(w, "  [ERROR] %s\n", e)
		}
	}

	if len(r.Conflicts) == 0 {
		fmt.Fprintln(w, "\nNo conflicts found.")
		return
	}

	fmt.Fprintf(w, "\nConflicts (%d):\n", len(r.Conflicts))
	for _, c := range r.Conflicts {
		fmt.Fprintf(w, "  [CONFLICT] line %d (%s) overlaps with line %d (%s)\n",
			c.A.Line, truncate(c.A.Command, 40),
			c.B.Line, truncate(c.B.Command, 40))
	}
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
