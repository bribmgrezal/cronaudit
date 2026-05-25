package audit

import (
	"fmt"
	"strings"

	"github.com/yourorg/cronaudit/internal/parser"
	"github.com/yourorg/cronaudit/internal/schedule"
)

// AnnotateReport holds per-entry annotations for a crontab file.
type AnnotateReport struct {
	File        string
	Annotations []schedule.Annotation
	ParseError  error
}

// AuditAnnotate parses a crontab file and annotates every schedule entry.
func AuditAnnotate(path string) AnnotateReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return AnnotateReport{File: path, ParseError: err}
	}

	var annotations []schedule.Annotation
	for _, e := range entries {
		a, err := schedule.Annotate(e.Schedule)
		if err != nil {
			// Store a minimal annotation recording the error as a warning.
			annotations = append(annotations, schedule.Annotation{
				Expression: e.Schedule,
				Warnings:   []string{err.Error()},
			})
			continue
		}
		annotations = append(annotations, a)
	}
	return AnnotateReport{File: path, Annotations: annotations}
}

// FormatAnnotateReport returns a human-readable summary of an AnnotateReport.
func FormatAnnotateReport(r AnnotateReport) string {
	if r.ParseError != nil {
		return fmt.Sprintf("[annotate] %s: parse error: %v\n", r.File, r.ParseError)
	}
	if len(r.Annotations) == 0 {
		return fmt.Sprintf("[annotate] %s: no entries\n", r.File)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[annotate] %s (%d entries)\n", r.File, len(r.Annotations)))
	for _, a := range r.Annotations {
		sb.WriteString(fmt.Sprintf("  expr:      %s\n", a.Expression))
		if a.Human != "" {
			sb.WriteString(fmt.Sprintf("  human:     %s\n", a.Human))
		}
		if a.Frequency != "" {
			sb.WriteString(fmt.Sprintf("  frequency: %s\n", a.Frequency))
		}
		if a.Risk != "" {
			sb.WriteString(fmt.Sprintf("  risk:      %s\n", a.Risk))
		}
		if len(a.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("  tags:      %s\n", strings.Join(a.Tags, ", ")))
		}
		if len(a.Suggestions) > 0 {
			sb.WriteString(fmt.Sprintf("  suggest:   %s\n", strings.Join(a.Suggestions, "; ")))
		}
		if len(a.Warnings) > 0 {
			sb.WriteString(fmt.Sprintf("  warnings:  %s\n", strings.Join(a.Warnings, "; ")))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
