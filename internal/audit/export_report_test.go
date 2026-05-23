package audit_test

import (
	"os"
	"strings"
	"testing"

	"github.com/cronaudit/internal/audit"
	"github.com/cronaudit/internal/schedule"
)

func writeTempExportCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "crontab-export-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAuditExport_JSONFormat(t *testing.T) {
	path := writeTempExportCrontab(t, "0 * * * * /usr/bin/hourly.sh\n*/5 * * * * /usr/bin/check.sh\n")
	r := audit.AuditExport(path, schedule.FormatJSON)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if !strings.Contains(r.Output, `"expression"`) {
		t.Error("expected JSON output with 'expression' field")
	}
	if !strings.Contains(r.Output, "hourly.sh") {
		t.Error("expected command in output")
	}
}

func TestAuditExport_TextFormat(t *testing.T) {
	path := writeTempExportCrontab(t, "30 6 * * 1-5 /usr/bin/weekday.sh\n")
	r := audit.AuditExport(path, schedule.FormatText)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if !strings.Contains(r.Output, "weekday.sh") {
		t.Error("expected command in text output")
	}
	if !strings.Contains(r.Output, "[line") {
		t.Error("expected line reference in text output")
	}
}

func TestAuditExport_ParseError(t *testing.T) {
	r := audit.AuditExport("/nonexistent/crontab", schedule.FormatJSON)
	if r.Error == "" {
		t.Error("expected error for missing file")
	}
}

func TestFormatExportReport_WithError(t *testing.T) {
	r := audit.ExportReport{
		File:   "test.crontab",
		Format: schedule.FormatJSON,
		Error:  "parse error: file not found",
	}
	out := audit.FormatExportReport(r)
	if !strings.Contains(out, "ERROR") {
		t.Error("expected ERROR in formatted output")
	}
}

func TestFormatExportReport_Success(t *testing.T) {
	r := audit.ExportReport{
		File:   "test.crontab",
		Format: schedule.FormatCSV,
		Output: "line,expression,command\n1,\"* * * * *\",\"/bin/true\"\n",
	}
	out := audit.FormatExportReport(r)
	if !strings.Contains(out, "Export Report") {
		t.Error("expected report header")
	}
	if !strings.Contains(out, "/bin/true") {
		t.Error("expected output content in formatted report")
	}
}
