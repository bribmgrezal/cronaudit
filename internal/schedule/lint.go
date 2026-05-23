package schedule

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning or suggestion for a cron expression.
type LintIssue struct {
	Field   string
	Message string
}

// LintResult holds all issues found for a given expression.
type LintResult struct {
	Expression string
	Issues     []LintIssue
}

// HasIssues returns true if any lint issues were found.
func (r LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

// Lint analyzes a cron expression for common mistakes and style issues.
// It does not validate correctness (use Validate for that).
func Lint(expr string) (LintResult, error) {
	result := LintResult{Expression: expr}

	fields, err := splitFieldsForTag(expr)
	if err != nil {
		return result, fmt.Errorf("lint: %w", err)
	}
	if len(fields) != 5 {
		return result, fmt.Errorf("lint: expected 5 fields, got %d", len(fields))
	}

	names := []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

	for i, field := range fields {
		issues := lintField(field, names[i])
		result.Issues = append(result.Issues, issues...)
	}

	// Warn about redundant DOM+DOW wildcard combo
	if fields[2] != "*" && fields[4] != "*" {
		result.Issues = append(result.Issues, LintIssue{
			Field:   "day-of-month+day-of-week",
			Message: "specifying both day-of-month and day-of-week may cause unexpected OR behaviour",
		})
	}

	return result, nil
}

func lintField(field, name string) []LintIssue {
	var issues []LintIssue

	// Detect leading zeros
	for _, part := range strings.Split(field, ",") {
		for _, token := range strings.FieldsFunc(part, func(r rune) bool { return r == '-' || r == '/' }) {
			if len(token) > 1 && token[0] == '0' && token != "0" {
				issues = append(issues, LintIssue{
					Field:   name,
					Message: fmt.Sprintf("leading zero in %q is non-standard; use Normalize to fix", token),
				})
			}
		}
	}

	// Detect */1 which is equivalent to *
	if field == "*/1" {
		issues = append(issues, LintIssue{
			Field:   name,
			Message: "*/1 is redundant; use * instead",
		})
	}

	// Detect range where start == end (e.g. 5-5)
	for _, part := range strings.Split(field, ",") {
		if idx := strings.Index(part, "-"); idx != -1 {
			base := part
			if s := strings.Index(part, "/"); s != -1 {
				base = part[:s]
			}
			sides := strings.SplitN(base, "-", 2)
			if len(sides) == 2 && sides[0] == sides[1] {
				issues = append(issues, LintIssue{
					Field:   name,
					Message: fmt.Sprintf("range %q has equal start and end; use a single value instead", base),
				})
			}
		}
	}

	return issues
}
