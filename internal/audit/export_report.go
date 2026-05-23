package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// ExportReport holds the result of exporting a crontab's schedules.
type ExportReport struct {
	File    string
	Format  schedule.ExportFormat
	Output  string
	Error   string
}

// AuditExport parses a crontab file and exports its entries in the given format.
func AuditExport(path string, format schedule.ExportFormat) ExportReport {
	report := ExportReport{File: path, Format: format}

	entries, err := parser.Parse(path)
	if err != nil {
		report.Error = fmt.Sprintf("parse error: %v", err)
		return report
	}

	var exports []schedule.ExportEntry
	for _, e := range entries {
		human := schedule.Humanize(e.Expression)
		tags := schedule.TagExpression(e.Expression)
		freq := "unknown"
		if f, err := schedule.ClassifyFrequency(e.Expression); err == nil {
			freq = string(f)
		}
		exports = append(exports, schedule.ExportEntry{
			Line:       e.LineNumber,
			Expression: e.Expression,
			Command:    e.Command,
			Human:      human,
			Tags:       tags,
			Frequency:  freq,
		})
	}

	out, err := schedule.ExportSchedules(exports, format)
	if err != nil {
		report.Error = fmt.Sprintf("export error: %v", err)
		return report
	}
	report.Output = out
	return report
}

// FormatExportReport returns a human-readable summary of an ExportReport.
func FormatExportReport(r ExportReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "=== Export Report: %s (format: %s) ===\n", r.File, r.Format)
	if r.Error != "" {
		fmt.Fprintf(&sb, "ERROR: %s\n", r.Error)
		return sb.String()
	}
	sb.WriteString(r.Output)
	return sb.String()
}
