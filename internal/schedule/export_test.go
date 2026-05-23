package schedule_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/internal/schedule"
)

var sampleEntries = []schedule.ExportEntry{
	{
		Line:       3,
		Expression: "0 * * * *",
		Command:    "/usr/bin/backup.sh",
		Human:      "every hour at minute 0",
		Tags:       []string{"hourly"},
		Frequency:  "hourly",
	},
	{
		Line:       7,
		Expression: "*/5 * * * *",
		Command:    "/usr/bin/check.sh",
		Human:      "every 5 minutes",
		Tags:       []string{"high-frequency"},
		Frequency:  "every-minute",
	},
}

func TestExportSchedules_JSON(t *testing.T) {
	out, err := schedule.ExportSchedules(sampleEntries, schedule.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"expression"`) {
		t.Error("expected JSON field 'expression' in output")
	}
	if !strings.Contains(out, "backup.sh") {
		t.Error("expected command in JSON output")
	}
}

func TestExportSchedules_CSV(t *testing.T) {
	out, err := schedule.ExportSchedules(sampleEntries, schedule.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "line,expression") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(out, "backup.sh") {
		t.Error("expected command in CSV output")
	}
}

func TestExportSchedules_Text(t *testing.T) {
	out, err := schedule.ExportSchedules(sampleEntries, schedule.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[line 3]") {
		t.Error("expected line number in text output")
	}
	if !strings.Contains(out, "every hour at minute 0") {
		t.Error("expected human description in text output")
	}
}

func TestExportSchedules_UnknownFormat(t *testing.T) {
	_, err := schedule.ExportSchedules(sampleEntries, schedule.ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExportSchedules_Empty(t *testing.T) {
	out, err := schedule.ExportSchedules(nil, schedule.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "null" && out != "[]" {
		// json.MarshalIndent on nil slice produces "null"
		if !strings.Contains(out, "null") && !strings.Contains(out, "[]") {
			t.Errorf("unexpected output for empty entries: %s", out)
		}
	}
}
