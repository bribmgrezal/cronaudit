package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// DependencyReport holds the result of a dependency audit.
type DependencyReport struct {
	File         string
	Dependencies []schedule.Dependency
	Error        string
}

// AuditDependencies parses a crontab file and finds implicit job dependencies
// within the given time window and gap threshold.
func AuditDependencies(path string, windowStart, windowEnd time.Time, gapMinutes int) DependencyReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return DependencyReport{File: path, Error: err.Error()}
	}

	type entry struct {
		Expr    string
		Command string
	}

	var inputs []struct {
		Expr    string
		Command string
	}
	for _, e := range entries {
		inputs = append(inputs, struct {
			Expr    string
			Command string
		}{Expr: e.Schedule, Command: e.Command})
	}

	deps, err := schedule.FindDependencies(inputs, windowStart, windowEnd, gapMinutes)
	if err != nil {
		return DependencyReport{File: path, Error: err.Error()}
	}
	return DependencyReport{File: path, Dependencies: deps}
}

// FormatDependencyReport returns a human-readable summary of the report.
func FormatDependencyReport(r DependencyReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Dependency Report: %s\n", r.File))
	if r.Error != "" {
		sb.WriteString(fmt.Sprintf("  error: %s\n", r.Error))
		return sb.String()
	}
	if len(r.Dependencies) == 0 {
		sb.WriteString("  no implicit dependencies detected\n")
		return sb.String()
	}
	for _, d := range r.Dependencies {
		sb.WriteString(fmt.Sprintf("  [leader]   %s (%s)\n", d.LeaderCommand, d.LeaderExpr))
		sb.WriteString(fmt.Sprintf("  [follower] %s (%s)\n", d.FollowerCommand, d.FollowerExpr))
		sb.WriteString(fmt.Sprintf("             gap: %d min — %s\n", d.GapMinutes, d.Note))
	}
	return sb.String()
}
