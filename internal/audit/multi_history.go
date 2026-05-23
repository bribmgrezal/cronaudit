package audit

import (
	"fmt"
	"time"
)

// HostHistoryResult pairs a host label with its HistoryReport.
type HostHistoryResult struct {
	Host   string
	Report HistoryReport
}

// AuditHostsHistory runs AuditHistory for each host file within [from, until)
// and returns per-host results.
func AuditHostsHistory(hosts map[string]string, from, until time.Time) []HostHistoryResult {
	results := make([]HostHistoryResult, 0, len(hosts))
	for host, path := range hosts {
		r := AuditHistory(path, from, until)
		results = append(results, HostHistoryResult{Host: host, Report: r})
	}
	return results
}

// FormatHistorySummary returns a multi-host summary of execution counts.
func FormatHistorySummary(results []HostHistoryResult) string {
	if len(results) == 0 {
		return "no hosts to summarise\n"
	}
	out := "History summary:\n"
	total := 0
	for _, hr := range results {
		if hr.Report.Err != nil {
			out += fmt.Sprintf("  %-20s  error: %v\n", hr.Host, hr.Report.Err)
			continue
		}
		n := len(hr.Report.Entries)
		total += n
		out += fmt.Sprintf("  %-20s  %d execution(s)\n", hr.Host, n)
	}
	out += fmt.Sprintf("total: %d execution(s) across %d host(s)\n", total, len(results))
	return out
}
