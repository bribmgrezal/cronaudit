package audit

import (
	"fmt"
	"strings"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// ClassifyEntry holds a crontab entry alongside its classification.
type ClassifyEntry struct {
	Line     int
	Command  string
	Expr     string
	Category schedule.Category
	Reason   string
}

// ClassifyReport is the result of classifying all entries in a crontab file.
type ClassifyReport struct {
	File    string
	Entries []ClassifyEntry
	Err     error
}

// AuditClassify parses a crontab file and classifies each entry's schedule.
func AuditClassify(path string) ClassifyReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return ClassifyReport{File: path, Err: err}
	}

	var results []ClassifyEntry
	for _, e := range entries {
		res, err := schedule.Classify(e.Schedule)
		if err != nil {
			res = schedule.ClassificationResult{
				Category: schedule.CategoryCustom,
				Reason:   fmt.Sprintf("parse error: %v", err),
			}
		}
		results = append(results, ClassifyEntry{
			Line:     e.Line,
			Command:  truncate(e.Command, 40),
			Expr:     e.Schedule,
			Category: res.Category,
			Reason:   res.Reason,
		})
	}
	return ClassifyReport{File: path, Entries: results}
}

// FormatClassifyReport returns a human-readable string for a ClassifyReport.
func FormatClassifyReport(r ClassifyReport) string {
	var sb strings.Builder
	if r.Err != nil {
		fmt.Fprintf(&sb, "[%s] error: %v\n", r.File, r.Err)
		return sb.String()
	}
	if len(r.Entries) == 0 {
		fmt.Fprintf(&sb, "[%s] no entries found\n", r.File)
		return sb.String()
	}
	fmt.Fprintf(&sb, "[%s] %d entries classified:\n", r.File, len(r.Entries))
	for _, e := range r.Entries {
		fmt.Fprintf(&sb, "  line %d: %-18s %-14s %s\n", e.Line, e.Expr, e.Category, e.Reason)
	}
	return sb.String()
}
