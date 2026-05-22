package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// MissedEntry describes a crontab entry that missed one or more executions.
type MissedEntry struct {
	Host    string
	Line    int
	Command string
	Expr    string
	Times   []time.Time
}

// MissedReport holds all missed-schedule findings across hosts.
type MissedReport struct {
	Window  schedule.MissedWindow
	Entries []MissedEntry
}

// AuditMissed checks each entry in entries for missed executions within window.
func AuditMissed(host string, entries []parser.Entry, window schedule.MissedWindow) (*MissedReport, error) {
	report := &MissedReport{Window: window}

	for _, entry := range entries {
		s, err := schedule.Parse(entry.Schedule)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", entry.Line, err)
		}

		times := schedule.FindMissed(s, window)
		if len(times) > 0 {
			report.Entries = append(report.Entries, MissedEntry{
				Host:    host,
				Line:    entry.Line,
				Command: entry.Command,
				Expr:    entry.Schedule,
				Times:   times,
			})
		}
	}

	return report, nil
}

// Format returns a human-readable summary of the missed report.
func (r *MissedReport) Format() string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("No missed schedules between %s and %s.",
			r.Window.Start.Format(time.RFC3339),
			r.Window.End.Format(time.RFC3339),
		)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Missed schedules (%s – %s):\n",
		r.Window.Start.Format(time.RFC3339),
		r.Window.End.Format(time.RFC3339),
	)
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  [%s] line %d: %q (%s) — %d missed\n",
			e.Host, e.Line, truncate(e.Command, 40), e.Expr, len(e.Times))
	}
	return sb.String()
}
