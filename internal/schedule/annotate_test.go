package schedule_test

import (
	"testing"

	"github.com/yourorg/cronaudit/internal/schedule"
)

func TestAnnotate_EveryMinute(t *testing.T) {
	a, err := schedule.Annotate("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Expression != "* * * * *" {
		t.Errorf("expression = %q, want %q", a.Expression, "* * * * *")
	}
	if a.Human == "" {
		t.Error("expected non-empty Human description")
	}
	if a.Frequency == "" {
		t.Error("expected non-empty Frequency")
	}
	if a.Risk == "" {
		t.Error("expected non-empty Risk level")
	}
	if len(a.Tags) == 0 {
		t.Error("expected at least one tag")
	}
}

func TestAnnotate_DailySchedule(t *testing.T) {
	a, err := schedule.Annotate("0 3 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Frequency == "" {
		t.Error("expected non-empty Frequency")
	}
	// daily at 03:00 should be low risk
	if a.Risk == "" {
		t.Error("expected non-empty Risk")
	}
}

func TestAnnotate_InvalidExpression(t *testing.T) {
	_, err := schedule.Annotate("99 * * * *")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestAnnotate_SuggestionsPopulated(t *testing.T) {
	// */1 should trigger a suggestion to use *
	a, err := schedule.Annotate("*/1 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.Suggestions) == 0 {
		t.Error("expected at least one suggestion for */1 expression")
	}
}

func TestAnnotate_WarningsPopulated(t *testing.T) {
	// */1 should also trigger a lint warning
	a, err := schedule.Annotate("*/1 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.Warnings) == 0 {
		t.Error("expected at least one lint warning for */1 expression")
	}
}

func TestAnnotate_WeekdaySchedule(t *testing.T) {
	a, err := schedule.Annotate("0 9 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.Tags) == 0 {
		t.Error("expected tags for weekday-restricted schedule")
	}
}
