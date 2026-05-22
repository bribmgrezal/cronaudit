package audit

import (
	"strings"
	"testing"
	"time"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

func TestAuditMissed_DetectsMissed(t *testing.T) {
	entries := []parser.Entry{
		{Line: 1, Schedule: "0 * * * *", Command: "/usr/bin/backup"},
	}

	start := time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	window := schedule.MissedWindow{Start: start, End: end}

	report, err := AuditMissed("host1", entries, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 missed entry, got %d", len(report.Entries))
	}
	if len(report.Entries[0].Times) != 3 {
		t.Errorf("expected 3 missed times, got %d", len(report.Entries[0].Times))
	}
}

func TestAuditMissed_NoMissed(t *testing.T) {
	entries := []parser.Entry{
		{Line: 1, Schedule: "0 3 * * *", Command: "/usr/bin/cleanup"},
	}

	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	window := schedule.MissedWindow{Start: start, End: end}

	report, err := AuditMissed("host1", entries, window)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected no missed entries, got %d", len(report.Entries))
	}
}

func TestAuditMissed_InvalidSchedule(t *testing.T) {
	entries := []parser.Entry{
		{Line: 5, Schedule: "bad schedule", Command: "/bin/foo"},
	}
	window := schedule.MissedWindow{
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}
	_, err := AuditMissed("host1", entries, window)
	if err == nil {
		t.Error("expected error for invalid schedule")
	}
}

func TestMissedReport_Format_Empty(t *testing.T) {
	window := schedule.MissedWindow{
		Start: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 15, 1, 0, 0, 0, time.UTC),
	}
	report := &MissedReport{Window: window}
	out := report.Format()
	if !strings.Contains(out, "No missed") {
		t.Errorf("expected 'No missed' in output, got: %s", out)
	}
}

func TestMissedReport_Format_WithEntries(t *testing.T) {
	window := schedule.MissedWindow{
		Start: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
	}
	report := &MissedReport{
		Window: window,
		Entries: []MissedEntry{
			{Host: "web1", Line: 3, Command: "/usr/bin/sync", Expr: "0 * * * *",
				Times: []time.Time{window.Start, window.End}},
		},
	}
	out := report.Format()
	if !strings.Contains(out, "web1") {
		t.Errorf("expected host name in output")
	}
	if !strings.Contains(out, "2 missed") {
		t.Errorf("expected missed count in output, got: %s", out)
	}
}
