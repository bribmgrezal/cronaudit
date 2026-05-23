package schedule

import (
	"testing"
)

func TestAssessRisk_LowRisk(t *testing.T) {
	r, err := AssessRisk("0 2 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Level != RiskLow {
		t.Errorf("expected low risk, got %s (reasons: %v)", r.Level, r.Reasons)
	}
}

func TestAssessRisk_HighFrequency(t *testing.T) {
	r, err := AssessRisk("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Level != RiskHigh {
		t.Errorf("expected high risk, got %s", r.Level)
	}
	if len(r.Reasons) == 0 {
		t.Error("expected at least one reason for high risk")
	}
}

func TestAssessRisk_WeekdayRestricted(t *testing.T) {
	r, err := AssessRisk("0 9 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Level == RiskLow {
		t.Errorf("expected medium or high risk for weekday-restricted, got low")
	}
}

func TestAssessRisk_InvalidExpression(t *testing.T) {
	_, err := AssessRisk("not a cron")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestAssessRisk_LintWarning(t *testing.T) {
	// */1 is a lint warning (redundant step of 1)
	r, err := AssessRisk("*/1 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// high frequency takes precedence, but reasons should include lint
	hasLint := false
	for _, reason := range r.Reasons {
		if len(reason) > 5 && reason[:4] == "lint" {
			hasLint = true
		}
	}
	if !hasLint {
		t.Errorf("expected lint reason in %v", r.Reasons)
	}
}
