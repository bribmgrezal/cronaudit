package audit

import (
	"os"
	"strings"
	"testing"
)

func writeTempRiskCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "crontab-risk-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAuditRisk_HighFrequency(t *testing.T) {
	path := writeTempRiskCrontab(t, "* * * * * root /usr/bin/poll\n")
	report := AuditRisk(path)
	if report.Error != "" {
		t.Fatalf("unexpected error: %s", report.Error)
	}
	if len(report.Entries) == 0 {
		t.Error("expected at least one risk entry for high-frequency schedule")
	}
}

func TestAuditRisk_LowRiskOnly(t *testing.T) {
	path := writeTempRiskCrontab(t, "0 3 * * * root /usr/bin/backup\n")
	report := AuditRisk(path)
	if report.Error != "" {
		t.Fatalf("unexpected error: %s", report.Error)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected no risk entries, got %d", len(report.Entries))
	}
}

func TestAuditRisk_ParseError(t *testing.T) {
	report := AuditRisk("/nonexistent/crontab.txt")
	if report.Error == "" {
		t.Error("expected parse error for missing file")
	}
}

func TestFormatRiskReport_WithEntries(t *testing.T) {
	path := writeTempRiskCrontab(t, "* * * * * root /bin/check\n")
	report := AuditRisk(path)
	out := FormatRiskReport(report)
	if !strings.Contains(out, "[risk]") {
		t.Errorf("expected '[risk]' in output, got: %s", out)
	}
}

func TestFormatRiskReport_NoIssues(t *testing.T) {
	report := RiskReport{File: "test.cron"}
	out := FormatRiskReport(report)
	if !strings.Contains(out, "no risk issues") {
		t.Errorf("expected 'no risk issues' in output, got: %s", out)
	}
}
