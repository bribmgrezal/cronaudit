package audit

import (
	"fmt"
	"strings"
	"time"
)

// HostSnapshotResult pairs a hostname with its SnapshotReport.
type HostSnapshotResult struct {
	Host   string
	Report SnapshotReport
}

// AuditHostsSnapshot runs AuditSnapshot across multiple hosts/files.
func AuditHostsSnapshot(hosts map[string]string, ref time.Time, n int) []HostSnapshotResult {
	results := make([]HostSnapshotResult, 0, len(hosts))
	for host, path := range hosts {
		r := AuditSnapshot(path, ref, n)
		results = append(results, HostSnapshotResult{Host: host, Report: r})
	}
	return results
}

// FormatSnapshotSummary renders a multi-host snapshot summary.
func FormatSnapshotSummary(results []HostSnapshotResult) string {
	var sb strings.Builder
	totalJobs := 0
	errorHosts := 0

	for _, hr := range results {
		if hr.Report.Error != "" {
			errorHosts++
			sb.WriteString(fmt.Sprintf("[%s] ERROR: %s\n", hr.Host, hr.Report.Error))
			continue
		}
		count := 0
		if hr.Report.Result != nil {
			count = len(hr.Report.Result.Entries)
		}
		totalJobs += count
		sb.WriteString(fmt.Sprintf("[%s] %d job(s) snapshotted\n", hr.Host, count))
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d host(s), %d total job(s), %d host(s) with errors\n",
		len(results), totalJobs, errorHosts))
	return sb.String()
}
