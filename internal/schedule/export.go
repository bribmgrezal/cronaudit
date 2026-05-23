package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExportFormat represents the output format for schedule exports.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// ExportEntry holds a single crontab entry's exported data.
type ExportEntry struct {
	Line       int      `json:"line"`
	Expression string   `json:"expression"`
	Command    string   `json:"command"`
	Human      string   `json:"human"`
	Tags       []string `json:"tags"`
	Frequency  string   `json:"frequency"`
}

// ExportSchedules converts a slice of ExportEntry to the requested format.
func ExportSchedules(entries []ExportEntry, format ExportFormat) (string, error) {
	switch format {
	case FormatJSON:
		return exportJSON(entries)
	case FormatCSV:
		return exportCSV(entries), nil
	case FormatText:
		return exportText(entries), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

func exportJSON(entries []ExportEntry) (string, error) {
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func exportCSV(entries []ExportEntry) string {
	var sb strings.Builder
	sb.WriteString("line,expression,command,human,frequency,tags\n")
	for _, e := range entries {
		fmt.Fprintf(&sb, "%d,%q,%q,%q,%s,%q\n",
			e.Line,
			e.Expression,
			e.Command,
			e.Human,
			e.Frequency,
			strings.Join(e.Tags, ";"),
		)
	}
	return sb.String()
}

func exportText(entries []ExportEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "[line %d] %s\n", e.Line, e.Expression)
		fmt.Fprintf(&sb, "  command  : %s\n", e.Command)
		fmt.Fprintf(&sb, "  human    : %s\n", e.Human)
		fmt.Fprintf(&sb, "  frequency: %s\n", e.Frequency)
		if len(e.Tags) > 0 {
			fmt.Fprintf(&sb, "  tags     : %s\n", strings.Join(e.Tags, ", "))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
