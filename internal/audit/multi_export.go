package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/schedule"
)

// HostExportResult holds the export result for a single host.
type HostExportResult struct {
	Host   string
	Report ExportReport
}

// AuditHostsExport runs AuditExport for each host crontab file.
func AuditHostsExport(hostFiles map[string]string, format schedule.ExportFormat) []HostExportResult {
	results := make([]HostExportResult, 0, len(hostFiles))
	for host, path := range hostFiles {
		r := AuditExport(path, format)
		results = append(results, HostExportResult{Host: host, Report: r})
	}
	return results
}

// FormatMultiExportSummary returns a combined export summary across all hosts.
func FormatMultiExportSummary(results []HostExportResult) string {
	var sb strings.Builder
	successCount := 0
	errorCount := 0

	for _, hr := range results {
		fmt.Fprintf(&sb, "--- Host: %s ---\n", hr.Host)
		if hr.Report.Error != "" {
			fmt.Fprintf(&sb, "  ERROR: %s\n", hr.Report.Error)
			errorCount++
		} else {
			lines := strings.Count(hr.Report.Output, "\n")
			fmt.Fprintf(&sb, "  exported %d lines (format: %s)\n", lines, hr.Report.Format)
			successCount++
		}
	}

	fmt.Fprintf(&sb, "\nSummary: %d host(s) exported successfully, %d error(s)\n",
		successCount, errorCount)
	return sb.String()
}
