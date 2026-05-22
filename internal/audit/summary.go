package audit

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// HostSummary holds aggregated audit results for a single host.
type HostSummary struct {
	Host        string
	Conflicts   int
	Missed      int
	TotalJobs   int
	AuditedAt   time.Time
	Errors      []string
}

// GlobalSummary aggregates HostSummary results across all audited hosts.
type GlobalSummary struct {
	Hosts       []HostSummary
	GeneratedAt time.Time
}

// TotalConflicts returns the sum of conflicts across all hosts.
func (g *GlobalSummary) TotalConflicts() int {
	total := 0
	for _, h := range g.Hosts {
		total += h.Conflicts
	}
	return total
}

// TotalMissed returns the sum of missed schedules across all hosts.
func (g *GlobalSummary) TotalMissed() int {
	total := 0
	for _, h := range g.Hosts {
		total += h.Missed
	}
	return total
}

// HostsWithIssues returns the count of hosts that have at least one conflict or missed job.
func (g *GlobalSummary) HostsWithIssues() int {
	count := 0
	for _, h := range g.Hosts {
		if h.Conflicts > 0 || h.Missed > 0 || len(h.Errors) > 0 {
			count++
		}
	}
	return count
}

// Format writes a human-readable summary report to w.
func (g *GlobalSummary) Format(w io.Writer) {
	fmt.Fprintf(w, "=== CronAudit Global Summary ===")
	fmt.Fprintf(w, "\nGenerated: %s\n", g.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Hosts audited : %d\n", len(g.Hosts))
	fmt.Fprintf(w, "Hosts with issues: %d\n", g.HostsWithIssues())
	fmt.Fprintf(w, "Total conflicts  : %d\n", g.TotalConflicts())
	fmt.Fprintf(w, "Total missed     : %d\n", g.TotalMissed())
	fmt.Fprintf(w, "\n%-30s %8s %8s %8s %s\n",
		"HOST", "JOBS", "CONFLICTS", "MISSED", "ERRORS")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 70))
	for _, h := range g.Hosts {
		errStr := "-"
		if len(h.Errors) > 0 {
			errStr = strings.Join(h.Errors, "; ")
			if len(errStr) > 30 {
				errStr = errStr[:27] + "..."
			}
		}
		fmt.Fprintf(w, "%-30s %8d %8d %8d %s\n",
			h.Host, h.TotalJobs, h.Conflicts, h.Missed, errStr)
	}
}
