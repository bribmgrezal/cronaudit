package audit

import (
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempDepCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "dep_crontab_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAuditDependencies_DetectsGap(t *testing.T) {
	path := writeTempDepCrontab(t, "0 * * * * leader.sh\n5 * * * * follower.sh\n")
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	r := AuditDependencies(path, start, end, 10)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if len(r.Dependencies) == 0 {
		t.Fatal("expected dependencies to be detected")
	}
}

func TestAuditDependencies_NoDeps(t *testing.T) {
	path := writeTempDepCrontab(t, "0 * * * * a.sh\n30 * * * * b.sh\n")
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	r := AuditDependencies(path, start, end, 10)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if len(r.Dependencies) != 0 {
		t.Errorf("expected no dependencies, got %d", len(r.Dependencies))
	}
}

func TestAuditDependencies_ParseError(t *testing.T) {
	r := AuditDependencies("/nonexistent/crontab.txt",
		time.Now(), time.Now().Add(time.Hour), 5)
	if r.Error == "" {
		t.Fatal("expected parse error")
	}
}

func TestFormatDependencyReport_WithEntries(t *testing.T) {
	r := DependencyReport{
		File: "test.crontab",
		Dependencies: []schedule.Dependency{
			{
				LeaderCommand:   "leader.sh",
				LeaderExpr:      "0 * * * *",
				FollowerCommand: "follower.sh",
				FollowerExpr:    "5 * * * *",
				GapMinutes:      5,
				Note:            "leader fires ~5 min before follower",
			},
		},
	}
	out := FormatDependencyReport(r)
	if !strings.Contains(out, "leader.sh") {
		t.Error("expected leader.sh in output")
	}
	if !strings.Contains(out, "follower.sh") {
		t.Error("expected follower.sh in output")
	}
	if !strings.Contains(out, "gap: 5") {
		t.Error("expected gap info in output")
	}
}

func TestFormatDependencyReport_Empty(t *testing.T) {
	r := DependencyReport{File: "empty.crontab"}
	out := FormatDependencyReport(r)
	if !strings.Contains(out, "no implicit dependencies") {
		t.Error("expected no-dependency message")
	}
}
