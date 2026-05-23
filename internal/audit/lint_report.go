package audit

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronaudit/internal/parser"
	"github.com/yourorg/cronaudit/internal/schedule"
)

// LintEntry holds lint results for a single crontab entry.
type LintEntry struct {
	Line       int
	Command    string
	Expression string
	Issues     []schedule.LintIssue
}

// LintReport summarises lint findings for a crontab file.
type LintReport struct {
	File    string
	Entries []LintEntry
	Err     error
}

// HasIssues returns true if any entry has lint issues.
func (r LintReport) HasIssues() bool {
	for _, e := range r.Entries {
		if len(e.Issues) > 0 {
			return true
		}
	}
	return false
}

// Format returns a human-readable summary of the lint report.
func (r LintReport) Format() string {
	if r.Err != nil {
		return fmt.Sprintf("lint report for %s: error: %v", r.File, r.Err)
	}
	if !r.HasIssues() {
		return fmt.Sprintf("lint report for %s: no issues found", r.File)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "lint report for %s:\n", r.File)
	for _, e := range r.Entries {
		if len(e.Issues) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "  line %d [%s] %s\n", e.Line, e.Expression, truncate(e.Command, 40))
		for _, issue := range e.Issues {
			fmt.Fprintf(&sb, "    [%s] %s\n", issue.Field, issue.Message)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// AuditLint parses a crontab file and runs lint checks on each entry.
func AuditLint(path string) LintReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return LintReport{File: path, Err: err}
	}

	report := LintReport{File: path}
	for _, entry := range entries {
		result, err := schedule.Lint(entry.Schedule)
		if err != nil {
			// skip unparseable expressions; validation covers these
			continue
		}
		report.Entries = append(report.Entries, LintEntry{
			Line:       entry.Line,
			Command:    entry.Command,
			Expression: entry.Schedule,
			Issues:     result.Issues,
		})
	}
	return report
}
