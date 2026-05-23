package audit

import (
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempSnapshotCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "crontab-snapshot-*.txt")
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

func TestAuditSnapshot_Basic(t *testing.T) {
	path := writeTempSnapshotCrontab(t, "* * * * * echo hello\n0 * * * * echo hourly\n")
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	r := AuditSnapshot(path, ref, 2)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(r.Result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.Result.Entries))
	}
	for _, e := range r.Result.Entries {
		if len(e.NextTimes) != 2 {
			t.Errorf("expected 2 next times for %q, got %d", e.Expression, len(e.NextTimes))
		}
	}
}

func TestAuditSnapshot_ParseError(t *testing.T) {
	r := AuditSnapshot("/nonexistent/crontab.txt", time.Now(), 3)
	if r.Error == "" {
		t.Error("expected error for missing file")
	}
}

func TestAuditSnapshot_InvalidExpression(t *testing.T) {
	path := writeTempSnapshotCrontab(t, "bad-expr cmd\n")
	r := AuditSnapshot(path, time.Now(), 2)
	if r.Error == "" {
		t.Error("expected error for invalid expression")
	}
}

func TestFormatSnapshotReport_WithEntries(t *testing.T) {
	path := writeTempSnapshotCrontab(t, "0 12 * * * /usr/bin/backup\n")
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	r := AuditSnapshot(path, ref, 2)
	out := FormatSnapshotReport(r)

	if !strings.Contains(out, "Snapshot:") {
		t.Error("expected 'Snapshot:' in output")
	}
	if !strings.Contains(out, "0 12 * * *") {
		t.Error("expected expression in output")
	}
	if !strings.Contains(out, "->") {
		t.Error("expected fire times with '->' in output")
	}
}

func TestFormatSnapshotReport_Error(t *testing.T) {
	r := SnapshotReport{File: "missing.txt", Error: "file not found"}
	out := FormatSnapshotReport(r)
	if !strings.Contains(out, "error") {
		t.Error("expected error in formatted output")
	}
}

func TestFormatSnapshotReport_Empty(t *testing.T) {
	path := writeTempSnapshotCrontab(t, "# just a comment\n")
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	r := AuditSnapshot(path, ref, 3)
	out := FormatSnapshotReport(r)
	if !strings.Contains(out, "no entries") {
		t.Error("expected 'no entries' for empty crontab")
	}
}
