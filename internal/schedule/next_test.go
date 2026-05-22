package schedule

import (
	"testing"
	"time"
)

func TestNext_HourlySchedule(t *testing.T) {
	s := mustParse(t, "0 * * * *")
	after := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	next, err := Next(s, after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNext_EveryMinute(t *testing.T) {
	s := mustParse(t, "* * * * *")
	after := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	next, err := Next(s, after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 31, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextN_ReturnsMultiple(t *testing.T) {
	s := mustParse(t, "0 9 * * 1") // every Monday at 09:00
	// Use a known Monday
	after := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Monday

	times, err := NextN(s, after, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 times, got %d", len(times))
	}

	// Each result should be a Monday at 09:00
	for _, ts := range times {
		if ts.Weekday() != time.Monday {
			t.Errorf("expected Monday, got %v", ts.Weekday())
		}
		if ts.Hour() != 9 || ts.Minute() != 0 {
			t.Errorf("expected 09:00, got %02d:%02d", ts.Hour(), ts.Minute())
		}
	}

	// Results should be 7 days apart
	gap := times[1].Sub(times[0])
	if gap != 7*24*time.Hour {
		t.Errorf("expected 7-day gap, got %v", gap)
	}
}

func TestNextN_InvalidN(t *testing.T) {
	s := mustParse(t, "* * * * *")
	after := time.Now()

	_, err := NextN(s, after, 0)
	if err == nil {
		t.Error("expected error for n=0")
	}

	_, err = NextN(s, after, 1001)
	if err == nil {
		t.Error("expected error for n=1001")
	}
}

func TestNext_AlreadyOnMinuteBoundary(t *testing.T) {
	s := mustParse(t, "* * * * *")
	// Exactly on a minute boundary — next should be one minute later
	after := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)

	next, err := Next(s, after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 3, 1, 12, 1, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}
