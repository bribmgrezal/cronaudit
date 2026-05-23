package schedule

import (
	"testing"
)

func TestNormalize_Alias(t *testing.T) {
	r, err := Normalize("@hourly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "0 * * * *" {
		t.Errorf("expected '0 * * * *', got %q", r.Normalized)
	}
	if len(r.Changes) == 0 {
		t.Error("expected at least one change recorded")
	}
}

func TestNormalize_WeekdayNames(t *testing.T) {
	r, err := Normalize("0 9 * * mon-fri")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "0 9 * * 1-5" {
		t.Errorf("expected '0 9 * * 1-5', got %q", r.Normalized)
	}
}

func TestNormalize_MonthNames(t *testing.T) {
	r, err := Normalize("0 0 1 jan *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "0 0 1 1 *" {
		t.Errorf("expected '0 0 1 1 *', got %q", r.Normalized)
	}
}

func TestNormalize_LeadingZeros(t *testing.T) {
	r, err := Normalize("05 09 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "5 9 * * *" {
		t.Errorf("expected '5 9 * * *', got %q", r.Normalized)
	}
	if len(r.Changes) == 0 {
		t.Error("expected changes for leading zeros")
	}
}

func TestNormalize_AlreadyCanonical(t *testing.T) {
	r, err := Normalize("0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "0 * * * *" {
		t.Errorf("expected unchanged, got %q", r.Normalized)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %v", r.Changes)
	}
}

func TestNormalize_InvalidFieldCount(t *testing.T) {
	_, err := Normalize("* * *")
	if err == nil {
		t.Error("expected error for wrong field count")
	}
}

func TestNormalize_AnnuallyAlias(t *testing.T) {
	r, err := Normalize("@annually")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Normalized != "0 0 1 1 *" {
		t.Errorf("expected '0 0 1 1 *', got %q", r.Normalized)
	}
}
