package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// RiskEntry holds risk assessment for a single crontab entry.
type RiskEntry struct {
	Line       int
	Command    string
	Expression string
	Result     schedule.RiskResult
}

// RiskReport summarises risk findings for a crontab file.
type RiskReport struct {
	File    string
	Entries []RiskEntry
	Error   string
}

// AuditRisk parses a crontab file and assesses risk for each entry.
func AuditRisk(path string) RiskReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return RiskReport{File: path, Error: err.Error()}
	}

	report := RiskReport{File: path}
	for _, e := range entries {
		result, err := schedule.AssessRisk(e.Schedule)
		if err != nil {
			continue
		}
		if result.Level != schedule.RiskLow {
			report.Entries = append(report.Entries, RiskEntry{
				Line:       e.Line,
				Command:    truncate(e.Command, 40),
				Expression: e.Schedule,
				Result:     result,
			})
		}
	}
	return report
}

// FormatRiskReport returns a human-readable summary of the risk report.
func FormatRiskReport(r RiskReport) string {
	if r.Error != "" {
		return fmt.Sprintf("[risk] %s: parse error: %s\n", r.File, r.Error)
	}
	if len(r.Entries) == 0 {
		return fmt.Sprintf("[risk] %s: no risk issues found\n", r.File)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "[risk] %s: %d issue(s)\n", r.File, len(r.Entries))
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  line %d [%s] %s — %s\n",
			e.Line, e.Result.Level, e.Expression, strings.Join(e.Result.Reasons, "; "))
	}
	return sb.String()
}
