package schedule

import (
	"testing"
	"time"
)

func mustParse(t *testing.T, expr string) *Schedule {
	t.Helper()
	s, err := Parse(expr)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", expr, err)
	}
	return s
}

func TestFindMissed_HourlyInWindow(t *testing.T) {
	s := mustParse(t, "0 * * * *") // every hour on the hour

	start := time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	missed := FindMissed(s, MissedWindow{Start: start, End: end})
	if len(missed) != 4 {
		t.Errorf("expected 4 missed events, got %d", len(missed))
	}
}

func TestFindMissed_NoMissed(t *testing.T) {
	s := mustParse(t, "30 2 * * *") // 02:30 daily

	start := time.Date(2024, 1, 15, 3, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 5, 0, 0, 0, time.UTC)

	missed := FindMissed(s, MissedWindow{Start: start, End: end})
	if len(missed) != 0 {
		t.Errorf("expected 0 missed events, got %d", len(missed))
	}
}

func TestMatches_ExactTime(t *testing.T) {
	s := mustParse(t, "30 14 * * 1") // 14:30 every Monday

	// 2024-01-15 is a Monday
	monday := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	if !Matches(s, monday) {
		t.Error("expected match on Monday 14:30")
	}

	tuesday := time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC)
	if Matches(s, tuesday) {
		t.Error("expected no match on Tuesday 14:30")
	}
}

func TestFindMissed_EveryMinute(t *testing.T) {
	s := mustParse(t, "* * * * *")

	start := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 0, 4, 0, 0, time.UTC)

	missed := FindMissed(s, MissedWindow{Start: start, End: end})
	if len(missed) != 5 {
		t.Errorf("expected 5 missed events, got %d", len(missed))
	}
}
