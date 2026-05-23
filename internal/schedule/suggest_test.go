package schedule

import (
	"testing"
)

func TestSuggest_HighFrequency(t *testing.T) {
	suggestions, err := Suggest("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion for every-minute schedule")
	}
	found := false
	for _, s := range suggestions {
		if s.Suggested == "*/5 * * * *" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected suggestion '*/5 * * * *', got: %+v", suggestions)
	}
}

func TestSuggest_WeekdayAlias(t *testing.T) {
	suggestions, err := Suggest("0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected suggestion for numeric weekday")
	}
	got := suggestions[0].Suggested
	want := "0 9 * * Mon"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSuggest_NoSuggestions(t *testing.T) {
	// Daily at midnight — no issues expected
	suggestions, err := Suggest("0 0 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions, got: %+v", suggestions)
	}
}

func TestSuggest_InvalidExpression(t *testing.T) {
	_, err := Suggest("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestSuggest_WeekdayWildcard_NoAlias(t *testing.T) {
	// Wildcard weekday should not trigger alias suggestion
	suggestions, err := Suggest("0 6 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, s := range suggestions {
		if s.Reason != "" && len(s.Reason) > 0 {
			if s.Original == "0 6 * * *" && s.Suggested != "0 6 * * *" {
				t.Errorf("unexpected suggestion for wildcard weekday: %+v", s)
			}
		}
	}
}

func TestSuggest_HighFrequencyWithWeekday(t *testing.T) {
	// Every minute on Monday — should get both suggestions
	suggestions, err := Suggest("* * * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) < 2 {
		t.Errorf("expected 2 suggestions, got %d: %+v", len(suggestions), suggestions)
	}
}
