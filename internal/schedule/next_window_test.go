package schedule

import (
	"testing"
	"time"
)

func makeWindow(start, end string) Window {
	loc := time.UTC
	parse := func(s string) time.Time {
		t, err := time.ParseInLocation("2006-01-02 15:04", s, loc)
		if err != nil {
			panic(err)
		}
		return t
	}
	return Window{Start: parse(start), End: parse(end)}
}

func TestNextInWindow_HourlySchedule(t *testing.T) {
	s, err := Parse("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	w := makeWindow("2024-01-01 00:00", "2024-01-01 05:00")
	times := NextInWindow(s, w)
	// Expect fires at 00:00, 01:00, 02:00, 03:00, 04:00, 05:00 => 6
	if len(times) != 6 {
		t.Errorf("expected 6 times, got %d: %v", len(times), times)
	}
}

func TestNextInWindow_EveryMinute(t *testing.T) {
	s, err := Parse("* * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	w := makeWindow("2024-01-01 00:00", "2024-01-01 00:04")
	times := NextInWindow(s, w)
	if len(times) != 5 {
		t.Errorf("expected 5 times, got %d", len(times))
	}
}

func TestNextInWindow_EmptyWindow(t *testing.T) {
	s, err := Parse("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	// End before Start
	w := Window{
		Start: time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 1, 4, 0, 0, 0, time.UTC),
	}
	times := NextInWindow(s, w)
	if len(times) != 0 {
		t.Errorf("expected 0 times for inverted window, got %d", len(times))
	}
}

func TestCountInWindow_Daily(t *testing.T) {
	s, err := Parse("0 9 * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	w := makeWindow("2024-01-01 00:00", "2024-01-07 23:59")
	count := CountInWindow(s, w)
	if count != 7 {
		t.Errorf("expected 7 occurrences over 7 days, got %d", count)
	}
}

func TestNextInWindow_NoFire(t *testing.T) {
	// Schedule fires on day 15 of each month at 00:00
	s, err := Parse("0 0 15 * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	// Window is entirely within Jan 1-14
	w := makeWindow("2024-01-01 00:00", "2024-01-14 23:59")
	times := NextInWindow(s, w)
	if len(times) != 0 {
		t.Errorf("expected 0 times, got %d: %v", len(times), times)
	}
}
