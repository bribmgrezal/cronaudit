package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// SnapshotReport holds the result of snapshotting a crontab file.
type SnapshotReport struct {
	File    string
	Error   string
	Result  *schedule.SnapshotResult
	N       int
}

// AuditSnapshot parses a crontab file and returns the next n fire times
// for each entry, relative to ref.
func AuditSnapshot(path string, ref time.Time, n int) SnapshotReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return SnapshotReport{File: path, Error: err.Error()}
	}

	var jobs [][2]string
	for _, e := range entries {
		jobs = append(jobs, [2]string{e.Schedule, e.Command})
	}

	result, err := schedule.Snapshot(jobs, ref, n)
	if err != nil {
		return SnapshotReport{File: path, Error: err.Error()}
	}

	return SnapshotReport{File: path, Result: result, N: n}
}

// FormatSnapshotReport renders a SnapshotReport as a human-readable string.
func FormatSnapshotReport(r SnapshotReport) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Snapshot: %s\n", r.File))

	if r.Error != "" {
		sb.WriteString(fmt.Sprintf("  error: %s\n", r.Error))
		return sb.String()
	}

	if r.Result == nil || len(r.Result.Entries) == 0 {
		sb.WriteString("  no entries\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("  generated at: %s\n", r.Result.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("  showing next %d fire time(s) per job\n", r.N))

	for _, e := range r.Result.Entries {
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", e.Expression, truncate(e.Command, 40)))
		for _, t := range e.NextTimes {
			sb.WriteString(fmt.Sprintf("      -> %s\n", t.Format(time.RFC3339)))
		}
	}
	return sb.String()
}
