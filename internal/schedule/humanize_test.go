package schedule

import (
	"testing"
)

func TestHumanize_EveryMinute(t *testing.T) {
	res, err := Humanize("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Description != "every minute" {
		t.Errorf("got %q, want %q", res.Description, "every minute")
	}
	if res.Frequency != "every_minute" {
		t.Errorf("got frequency %q, want %q", res.Frequency, "every_minute")
	}
}

func TestHumanize_HourlyAtMinute(t *testing.T) {
	res, err := Humanize("30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "at minute 30 of every hour"
	if res.Description != want {
		t.Errorf("got %q, want %q", res.Description, want)
	}
}

func TestHumanize_DailyAtTime(t *testing.T) {
	res, err := Humanize("0 9 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "at 09:00 every day"
	if res.Description != want {
		t.Errorf("got %q, want %q", res.Description, want)
	}
}

func TestHumanize_StepMinutes(t *testing.T) {
	res, err := Humanize("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "every 5 minute(s)"
	if res.Description != want {
		t.Errorf("got %q, want %q", res.Description, want)
	}
}

func TestHumanize_WeekdaySchedule(t *testing.T) {
	res, err := Humanize("0 8 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "at 08:00 on weekday 1"
	if res.Description != want {
		t.Errorf("got %q, want %q", res.Description, want)
	}
}

func TestHumanize_MonthDay(t *testing.T) {
	res, err := Humanize("0 6 1 1 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "at 06:00 on 1/1"
	if res.Description != want {
		t.Errorf("got %q, want %q", res.Description, want)
	}
}

func TestHumanize_InvalidExpression(t *testing.T) {
	_, err := Humanize("99 * * * *")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestHumanize_ExpressionPreserved(t *testing.T) {
	expr := "15 14 * * *"
	res, err := Humanize(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expression != expr {
		t.Errorf("expression not preserved: got %q, want %q", res.Expression, expr)
	}
}
