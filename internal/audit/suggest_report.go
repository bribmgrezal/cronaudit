package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/cronaudit/internal/parser"
	"github.com/cronaudit/cronaudit/internal/schedule"
)

// SuggestReport holds improvement suggestions for a crontab file.
type SuggestReport struct {
	File        string
	Suggestions []EntrySuggestion
	ParseError  error
}

// EntrySuggestion pairs a crontab entry with its schedule suggestions.
type EntrySuggestion struct {
	Line        int
	Command     string
	Suggestions []schedule.Suggestion
}

// AuditSuggestions parses a crontab file and returns schedule improvement suggestions.
func AuditSuggestions(path string) SuggestReport {
	entries, err := parser.Parse(path)
	if err != nil {
		return SuggestReport{File: path, ParseError: err}
	}

	var results []EntrySuggestion
	for _, entry := range entries {
		suggestions, err := schedule.Suggest(entry.Schedule)
		if err != nil || len(suggestions) == 0 {
			continue
		}
		results = append(results, EntrySuggestion{
			Line:        entry.Line,
			Command:     entry.Command,
			Suggestions: suggestions,
		})
	}

	return SuggestReport{
		File:        path,
		Suggestions: results,
	}
}

// Format returns a human-readable summary of the suggestion report.
func (r SuggestReport) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Suggestions for: %s\n", r.File))

	if r.ParseError != nil {
		sb.WriteString(fmt.Sprintf("  parse error: %v\n", r.ParseError))
		return sb.String()
	}

	if len(r.Suggestions) == 0 {
		sb.WriteString("  no suggestions\n")
		return sb.String()
	}

	for _, es := range r.Suggestions {
		sb.WriteString(fmt.Sprintf("  line %d: %s\n", es.Line, truncate(es.Command, 40)))
		for _, s := range es.Suggestions {
			sb.WriteString(fmt.Sprintf("    - %s\n", s.Reason))
			sb.WriteString(fmt.Sprintf("      try: %s\n", s.Suggested))
		}
	}

	return sb.String()
}
