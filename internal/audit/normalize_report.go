package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// NormalizeEntry holds the result of normalizing a single crontab entry.
type NormalizeEntry struct {
	Line       int
	Command    string
	Original   string
	Normalized string
	Changes    []string
}

// NormalizeReport holds the full normalization audit results.
type NormalizeReport struct {
	File    string
	Entries []NormalizeEntry
	Errors  []string
}

// AuditNormalize parses a crontab file and normalizes each schedule expression.
func AuditNormalize(path string) (NormalizeReport, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return NormalizeReport{File: path}, fmt.Errorf("parse error: %w", err)
	}

	report := NormalizeReport{File: path}
	for _, e := range entries {
		result, nerr := schedule.Normalize(e.Schedule)
		if nerr != nil {
			report.Errors = append(report.Errors,
				fmt.Sprintf("line %d: %v", e.Line, nerr))
			continue
		}
		report.Entries = append(report.Entries, NormalizeEntry{
			Line:       e.Line,
			Command:    e.Command,
			Original:   result.Original,
			Normalized: result.Normalized,
			Changes:    result.Changes,
		})
	}
	return report, nil
}

// Format returns a human-readable summary of the normalization report.
func (r NormalizeReport) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Normalize Report: %s\n", r.File))

	if len(r.Errors) > 0 {
		sb.WriteString("Errors:\n")
		for _, e := range r.Errors {
			sb.WriteString(fmt.Sprintf("  %s\n", e))
		}
	}

	changed := 0
	for _, e := range r.Entries {
		if len(e.Changes) > 0 {
			changed++
			sb.WriteString(fmt.Sprintf("  line %d [%s]: %q -> %q\n",
				e.Line, truncate(e.Command, 30), e.Original, e.Normalized))
			for _, c := range e.Changes {
				sb.WriteString(fmt.Sprintf("    - %s\n", c))
			}
		}
	}

	if changed == 0 && len(r.Errors) == 0 {
		sb.WriteString("  All expressions are already canonical.\n")
	}
	return sb.String()
}
