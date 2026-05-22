package audit

import (
	"fmt"
	"strings"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// TaggedEntry pairs a crontab entry with its derived tags.
type TaggedEntry struct {
	Line    int
	Command string
	Expr    string
	Tags    []schedule.Tag
}

// TagReport holds the result of tagging all entries in a crontab file.
type TagReport struct {
	File    string
	Entries []TaggedEntry
	Errors  []string
}

// AuditTags parses a crontab file and tags each entry's schedule expression.
func AuditTags(filename, content string) TagReport {
	report := TagReport{File: filename}

	entries, err := parser.Parse(content)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("parse error: %v", err))
		return report
	}

	for _, e := range entries {
		tags, err := schedule.TagExpression(e.Schedule)
		if err != nil {
			report.Errors = append(report.Errors,
				fmt.Sprintf("line %d: %v", e.Line, err))
			continue
		}
		report.Entries = append(report.Entries, TaggedEntry{
			Line:    e.Line,
			Command: e.Command,
			Expr:    e.Schedule,
			Tags:    tags,
		})
	}

	return report
}

// Format returns a human-readable summary of the tag report.
func (r TagReport) Format() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Tag report for %s\n", r.File)

	if len(r.Errors) > 0 {
		fmt.Fprintln(&sb, "Errors:")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "  - %s\n", e)
		}
	}

	if len(r.Entries) == 0 {
		fmt.Fprintln(&sb, "  (no tagged entries)")
		return sb.String()
	}

	for _, te := range r.Entries {
		keys := make([]string, len(te.Tags))
		for i, tg := range te.Tags {
			keys[i] = tg.Key
		}
		fmt.Fprintf(&sb, "  line %d [%s] %s => [%s]\n",
			te.Line, te.Expr, truncate(te.Command, 40), strings.Join(keys, ", "))
	}

	return sb.String()
}
