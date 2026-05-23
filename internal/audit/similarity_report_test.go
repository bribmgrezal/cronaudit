package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCrontab(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "crontab")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp crontab: %v", err)
	}
	return p
}

func TestAuditSimilarity_DetectsGroup(t *testing.T) {
	content := "0 * * * * /usr/bin/job1\n0 * * * * /usr/bin/job2\n"
	p := writeTempCrontab(t, content)

	report, err := AuditSimilarity(p, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(report.Groups))
	}
	if len(report.Groups[0].Entries) != 2 {
		t.Errorf("expected 2 entries in group, got %d", len(report.Groups[0].Entries))
	}
}

func TestAuditSimilarity_NoGroups(t *testing.T) {
	content := "0 * * * * /usr/bin/job1\n30 * * * * /usr/bin/job2\n"
	p := writeTempCrontab(t, content)

	report, err := AuditSimilarity(p, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(report.Groups))
	}
}

func TestAuditSimilarity_ParseError(t *testing.T) {
	_, err := AuditSimilarity("/nonexistent/crontab", 0.8)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSimilarityReport_Format_Empty(t *testing.T) {
	r := &SimilarityReport{File: "test", Groups: nil}
	out := r.Format()
	if !strings.Contains(out, "No similar") {
		t.Errorf("expected 'No similar' in output, got: %s", out)
	}
}

func TestSimilarityReport_Format_WithGroups(t *testing.T) {
	content := "0 * * * * /usr/bin/a\n0 * * * * /usr/bin/b\n"
	p := writeTempCrontab(t, content)

	report, err := AuditSimilarity(p, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := report.Format()
	if !strings.Contains(out, "Group 1") {
		t.Errorf("expected 'Group 1' in output, got: %s", out)
	}
	if !strings.Contains(out, "/usr/bin/a") {
		t.Errorf("expected command in output, got: %s", out)
	}
}
