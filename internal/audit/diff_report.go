package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/cronaudit/internal/schedule"
)

// DiffReport holds the result of comparing two cron expressions.
type DiffReport struct {
	ExprA   string
	ExprB   string
	From    time.Time
	To      time.Time
	Result  *schedule.DiffResult
	Err     error
}

// AuditDiff compares two cron expressions over the given time window.
func AuditDiff(exprA, exprB string, from, to time.Time) *DiffReport {
	r := &DiffReport{
		ExprA: exprA,
		ExprB: exprB,
		From:  from,
		To:    to,
	}
	result, err := schedule.DiffSchedules(exprA, exprB, from, to)
	if err != nil {
		r.Err = err
		return r
	}
	r.Result = result
	return r
}

// Format returns a human-readable summary of the diff report.
func (r *DiffReport) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Diff: [%s] vs [%s]\n", r.ExprA, r.ExprB))
	sb.WriteString(fmt.Sprintf("Window: %s — %s\n",
		r.From.Format(time.RFC3339), r.To.Format(time.RFC3339)))

	if r.Err != nil {
		sb.WriteString(fmt.Sprintf("Error: %v\n", r.Err))
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("Common firings:    %d\n", len(r.Result.Common)))
	sb.WriteString(fmt.Sprintf("Only in A (%s): %d\n", r.ExprA, len(r.Result.OnlyInA)))
	sb.WriteString(fmt.Sprintf("Only in B (%s): %d\n", r.ExprB, len(r.Result.OnlyInB)))
	return sb.String()
}
