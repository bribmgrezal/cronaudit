package schedule

import (
	"strings"
	"testing"
	"time"
)

func makeTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestBuildHistory_HourlyInWindow(t *testing.T) {
	from := makeTime("2024-01-01 00:00")
	until := makeTime("2024-01-01 03:00")
	entries, err := BuildHistory("0 * * * *", "/bin/hourly", from, until)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Time != makeTime("2024-01-01 00:00") {
		t.Errorf("unexpected first entry: %v", entries[0].Time)
	}
}

func TestBuildHistory_NoFirings(t *testing.T) {
	// daily at midnight; window is 01:00–02:00 so no match
	from := makeTime("2024-01-01 01:00")
	until := makeTime("2024-01-01 02:00")
	entries, err := BuildHistory("0 0 * * *", "/bin/daily", from, until)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestBuildHistory_InvalidExpression(t *testing.T) {
	from := makeTime("2024-01-01 00:00")
	until := makeTime("2024-01-01 01:00")
	_, err := BuildHistory("not-a-cron", "/bin/x", from, until)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestFormatHistory_Empty(t *testing.T) {
	out := FormatHistory(nil)
	if !strings.Contains(out, "no executions") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatHistory_WithEntries(t *testing.T) {
	entries := []HistoryEntry{
		{Time: makeTime("2024-01-01 00:00"), Expression: "0 * * * *", Command: "/bin/hourly"},
		{Time: makeTime("2024-01-01 01:00"), Expression: "0 * * * *", Command: "/bin/hourly"},
	}
	out := FormatHistory(entries)
	if !strings.Contains(out, "2 execution") {
		t.Errorf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "/bin/hourly") {
		t.Errorf("expected command in output, got: %q", out)
	}
}
