package audit

import (
	"fmt"
	"strings"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// ValidationIssue pairs a crontab entry with the errors found in its schedule.
type ValidationIssue struct {
	Line    int
	Command string
	Errors  []error
}

// ValidationReport holds the results of validating all entries in a crontab.
type ValidationReport struct {
	Host   string
	Issues []ValidationIssue
}

// Valid returns true when no validation issues were found.
func (r ValidationReport) Valid() bool { return len(r.Issues) == 0 }

// Format returns a human-readable summary of the report.
func (r ValidationReport) Format() string {
	if r.Valid() {
		return fmt.Sprintf("[%s] all schedules valid", r.Host)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] %d invalid schedule(s):\n", r.Host, len(r.Issues))
	for _, issue := range r.Issues {
		fmt.Fprintf(&sb, "  line %d (%s):\n", issue.Line, truncate(issue.Command, 40))
		for _, e := range issue.Errors {
			fmt.Fprintf(&sb, "    - %s\n", e.Error())
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// AuditValidation parses the crontab content and validates every entry's
// schedule expression, returning a ValidationReport.
func AuditValidation(host, content string) ValidationReport {
	report := ValidationReport{Host: host}
	entries, err := parser.Parse(content)
	if err != nil {
		report.Issues = append(report.Issues, ValidationIssue{
			Command: "<parse error>",
			Errors:  []error{err},
		})
		return report
	}
	for _, e := range entries {
		expr := strings.Join([]string{
			e.Minute, e.Hour, e.DayOfMonth, e.Month, e.DayOfWeek,
		}, " ")
		if errs := schedule.Validate(expr); len(errs) > 0 {
			report.Issues = append(report.Issues, ValidationIssue{
				Line:    e.Line,
				Command: e.Command,
				Errors:  errs,
			})
		}
	}
	return report
}
