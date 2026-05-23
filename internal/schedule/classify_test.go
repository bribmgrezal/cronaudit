package schedule

import (
	"testing"
)

func TestClassify_EveryMinute(t *testing.T) {
	res, err := Classify("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryEveryMinute {
		t.Errorf("expected %s, got %s", CategoryEveryMinute, res.Category)
	}
}

func TestClassify_HighFrequency(t *testing.T) {
	res, err := Classify("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryHighFreq {
		t.Errorf("expected %s, got %s", CategoryHighFreq, res.Category)
	}
}

func TestClassify_Hourly(t *testing.T) {
	res, err := Classify("30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryHourly {
		t.Errorf("expected %s, got %s", CategoryHourly, res.Category)
	}
}

func TestClassify_Daily(t *testing.T) {
	res, err := Classify("0 6 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryDaily {
		t.Errorf("expected %s, got %s", CategoryDaily, res.Category)
	}
}

func TestClassify_Weekly(t *testing.T) {
	res, err := Classify("0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryWeekly {
		t.Errorf("expected %s, got %s", CategoryWeekly, res.Category)
	}
}

func TestClassify_Monthly(t *testing.T) {
	res, err := Classify("0 0 1 * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryMonthly {
		t.Errorf("expected %s, got %s", CategoryMonthly, res.Category)
	}
}

func TestClassify_Custom(t *testing.T) {
	res, err := Classify("0 */3 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Category != CategoryCustom {
		t.Errorf("expected %s, got %s", CategoryCustom, res.Category)
	}
}

func TestClassify_InvalidExpression(t *testing.T) {
	_, err := Classify("bad expression")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestClassify_ReasonNotEmpty(t *testing.T) {
	cases := []string{
		"* * * * *",
		"*/5 * * * *",
		"30 * * * *",
		"0 6 * * *",
		"0 9 * * 1",
		"0 0 1 * *",
		"0 */3 * * 1-5",
	}
	for _, expr := range cases {
		res, err := Classify(expr)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", expr, err)
			continue
		}
		if res.Reason == "" {
			t.Errorf("%q: expected non-empty reason", expr)
		}
	}
}
