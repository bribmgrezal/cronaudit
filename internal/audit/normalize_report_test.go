package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempNormalizeCrontab(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "crontab")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp crontab: %v", err)
	}
	return p
}

func TestAuditNormalize_AlreadyCanonical(t *testing.T) {
	p := writeTempNormalizeCrontab(t, "0 * * * * /usr/bin/backup\n")
	report, err := AuditNormalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Errors) != 0 {
		t.Errorf("expected no errors, got %v", report.Errors)
	}
	for _, e := range report.Entries {
		if len(e.Changes) > 0 {
			t.Errorf("expected no changes for canonical entry, got %v", e.Changes)
		}
	}
}

func TestAuditNormalize_DetectsChanges(t *testing.T) {
	p := writeTempNormalizeCrontab(t, "05 09 * * mon /usr/bin/report\n")
	report, err := AuditNormalize(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	e := report.Entries[0]
	if e.Normalized == e.Original {
		t.Error("expected normalization to change the expression")
	}
	if len(e.Changes) == 0 {
		t.Error("expected changes to be recorded")
	}
}

func TestAuditNormalize_ParseError(t *testing.T) {
	_, err := AuditNormalize("/nonexistent/path/crontab")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestNormalizeReport_Format_AllCanonical(t *testing.T) {
	p := writeTempNormalizeCrontab(t, "0 0 * * * /bin/daily\n")
	report, _ := AuditNormalize(p)
	out := report.Format()
	if !strings.Contains(out, "canonical") {
		t.Errorf("expected 'canonical' in output, got: %s", out)
	}
}

func TestNormalizeReport_Format_WithChanges(t *testing.T) {
	p := writeTempNormalizeCrontab(t, "@daily /bin/cleanup\n")
	report, _ := AuditNormalize(p)
	out := report.Format()
	if !strings.Contains(out, "->") {
		t.Errorf("expected '->' in output for changed entry, got: %s", out)
	}
}
