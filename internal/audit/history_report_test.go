package audit

import (
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempHistoryCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "crontab-hist-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAuditHistory_DetectsEntries(t *testing.T) {
	path := writeTempHistoryCrontab(t, "0 * * * * /bin/hourly\n")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC)
	r := AuditHistory(path, from, until)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestAuditHistory_NoEntries(t *testing.T) {
	path := writeTempHistoryCrontab(t, "0 0 * * * /bin/daily\n")
	from := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC)
	r := AuditHistory(path, from, until)
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestAuditHistory_ParseError(t *testing.T) {
	r := AuditHistory("/nonexistent/crontab", time.Now(), time.Now().Add(time.Hour))
	if r.Err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFormatHistoryReport_WithError(t *testing.T) {
	r := HistoryReport{File: "bad.crontab", Err: os.ErrNotExist}
	out := FormatHistoryReport(r)
	if !strings.Contains(out, "error") {
		t.Errorf("expected error text, got: %q", out)
	}
}

func TestFormatHistoryReport_WithEntries(t *testing.T) {
	path := writeTempHistoryCrontab(t, "0 * * * * /bin/hourly\n")
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC)
	r := AuditHistory(path, from, until)
	out := FormatHistoryReport(r)
	if !strings.Contains(out, "History for") {
		t.Errorf("expected header in output, got: %q", out)
	}
	if !strings.Contains(out, "/bin/hourly") {
		t.Errorf("expected command in output, got: %q", out)
	}
}
