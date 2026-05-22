package schedule

import (
	"testing"
)

func TestClassifyFrequency_EveryMinute(t *testing.T) {
	s, err := Parse("* * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	f := ClassifyFrequency(s)
	if f.Label != "every minute" {
		t.Errorf("expected 'every minute', got %q", f.Label)
	}
	if f.EstimatedPerDay != 1440 {
		t.Errorf("expected 1440 per day, got %v", f.EstimatedPerDay)
	}
}

func TestClassifyFrequency_Hourly(t *testing.T) {
	s, err := Parse("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	f := ClassifyFrequency(s)
	if f.Label != "hourly" {
		t.Errorf("expected 'hourly', got %q", f.Label)
	}
	if f.EstimatedPerDay != 24 {
		t.Errorf("expected 24 per day, got %v", f.EstimatedPerDay)
	}
}

func TestClassifyFrequency_Daily(t *testing.T) {
	s, err := Parse("0 9 * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	f := ClassifyFrequency(s)
	if f.Label != "daily" {
		t.Errorf("expected 'daily', got %q", f.Label)
	}
	if f.EstimatedPerDay != 1 {
		t.Errorf("expected 1 per day, got %v", f.EstimatedPerDay)
	}
}

func TestClassifyFrequency_TwiceDaily(t *testing.T) {
	s, err := Parse("0 9,18 * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	f := ClassifyFrequency(s)
	if f.Label != "twice daily" {
		t.Errorf("expected 'twice daily', got %q", f.Label)
	}
}

func TestIsHighFrequency_EveryMinute(t *testing.T) {
	s, err := Parse("* * * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !IsHighFrequency(s) {
		t.Error("expected every-minute schedule to be high frequency")
	}
}

func TestIsHighFrequency_Daily(t *testing.T) {
	s, err := Parse("0 9 * * *")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if IsHighFrequency(s) {
		t.Error("expected daily schedule to not be high frequency")
	}
}
