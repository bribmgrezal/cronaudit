package schedule

import (
	"testing"
	"time"
)

func TestDiffSchedules_NoOverlap(t *testing.T) {
	// A fires at :00, B fires at :30 — no common minutes in a 1-hour window
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 59, 0, 0, time.UTC)

	result, err := DiffSchedules("0 * * * *", "30 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Common) != 0 {
		t.Errorf("expected 0 common, got %d", len(result.Common))
	}
	if len(result.OnlyInA) != 1 {
		t.Errorf("expected 1 in A, got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 1 {
		t.Errorf("expected 1 in B, got %d", len(result.OnlyInB))
	}
}

func TestDiffSchedules_FullOverlap(t *testing.T) {
	// Both expressions are identical
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC)

	result, err := DiffSchedules("0 * * * *", "0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.OnlyInA) != 0 {
		t.Errorf("expected 0 only in A, got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 0 {
		t.Errorf("expected 0 only in B, got %d", len(result.OnlyInB))
	}
	if len(result.Common) != 3 {
		t.Errorf("expected 3 common, got %d", len(result.Common))
	}
}

func TestDiffSchedules_PartialOverlap(t *testing.T) {
	// A: every 30 min (0,30), B: every hour at :00 — :00 is common
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 59, 0, 0, time.UTC)

	result, err := DiffSchedules("*/30 * * * *", "0 * * * *", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Common) != 1 {
		t.Errorf("expected 1 common, got %d", len(result.Common))
	}
	if len(result.OnlyInA) != 1 {
		t.Errorf("expected 1 only in A (the :30), got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 0 {
		t.Errorf("expected 0 only in B, got %d", len(result.OnlyInB))
	}
}

func TestDiffSchedules_InvalidExprA(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	_, err := DiffSchedules("bad expr", "0 * * * *", from, to)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestDiffSchedules_InvalidExprB(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	_, err := DiffSchedules("0 * * * *", "bad expr", from, to)
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}
