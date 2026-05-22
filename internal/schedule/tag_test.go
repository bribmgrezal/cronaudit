package schedule

import (
	"testing"
)

func TestTagExpression_EveryMinute(t *testing.T) {
	tags, err := TagExpression("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertHasTag(t, tags, "frequency")
	assertHasTag(t, tags, "high-frequency")
}

func TestTagExpression_Hourly(t *testing.T) {
	tags, err := TagExpression("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertHasTag(t, tags, "frequency")
	assertHasTag(t, tags, "fixed-minute")
	assertNoTag(t, tags, "high-frequency")
}

func TestTagExpression_WeekdayRestricted(t *testing.T) {
	tags, err := TagExpression("0 9 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertHasTag(t, tags, "weekday-restricted")
	assertHasTag(t, tags, "fixed-minute")
}

func TestTagExpression_MonthRestricted(t *testing.T) {
	tags, err := TagExpression("0 0 1 6 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertHasTag(t, tags, "month-restricted")
}

func TestTagExpression_Invalid(t *testing.T) {
	_, err := TagExpression("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestTagExpression_StepNoFixedMinute(t *testing.T) {
	tags, err := TagExpression("*/15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNoTag(t, tags, "fixed-minute")
}

// helpers

func assertHasTag(t *testing.T, tags []Tag, key string) {
	t.Helper()
	for _, tg := range tags {
		if tg.Key == key {
			return
		}
	}
	t.Errorf("expected tag %q not found in %v", key, tags)
}

func assertNoTag(t *testing.T, tags []Tag, key string) {
	t.Helper()
	for _, tg := range tags {
		if tg.Key == key {
			t.Errorf("unexpected tag %q found in %v", key, tags)
			return
		}
	}
}
