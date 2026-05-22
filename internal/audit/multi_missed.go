package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// HostCrontab pairs a host name with its raw crontab content.
type HostCrontab struct {
	Host    string
	Content string
}

// AuditHostsMissed parses crontabs for multiple hosts and aggregates missed
// schedule findings within the given window.
func AuditHostsMissed(hosts []HostCrontab, window schedule.MissedWindow) ([]*MissedReport, []error) {
	var reports []*MissedReport
	var errs []error

	for _, h := range hosts {
		entries, err := parser.Parse(strings.NewReader(h.Content))
		if err != nil {
			errs = append(errs, fmt.Errorf("host %s: parse error: %w", h.Host, err))
			continue
		}

		report, err := AuditMissed(h.Host, entries, window)
		if err != nil {
			errs = append(errs, fmt.Errorf("host %s: audit error: %w", h.Host, err))
			continue
		}

		reports = append(reports, report)
	}

	return reports, errs
}

// SummaryMissed returns a combined text summary of all missed reports.
func SummaryMissed(reports []*MissedReport) string {
	var sb strings.Builder
	total := 0
	for _, r := range reports {
		total += len(r.Entries)
		sb.WriteString(r.Format())
		sb.WriteString("\n")
	}
	if total == 0 {
		return "All schedules accounted for — no missed executions detected.\n"
	}
	fmt.Fprintf(&sb, "Total hosts with missed schedules: %d\n", countHostsWithMissed(reports))
	return sb.String()
}

func countHostsWithMissed(reports []*MissedReport) int {
	count := 0
	for _, r := range reports {
		if len(r.Entries) > 0 {
			count++
		}
	}
	return count
}
