package schedule

import (
	"testing"
	"time"
)

func TestFindDependencies_DetectsGap(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	entries := []struct {
		Expr    string
		Command string
	}{
		{"0 * * * *", "leader.sh"},
		{"5 * * * *", "follower.sh"},
	}

	deps, err := FindDependencies(entries, start, end, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) == 0 {
		t.Fatal("expected at least one dependency")
	}
	d := deps[0]
	if d.LeaderCommand != "leader.sh" {
		t.Errorf("expected leader leader.sh, got %s", d.LeaderCommand)
	}
	if d.FollowerCommand != "follower.sh" {
		t.Errorf("expected follower follower.sh, got %s", d.FollowerCommand)
	}
	if d.GapMinutes != 5 {
		t.Errorf("expected gap 5, got %d", d.GapMinutes)
	}
}

func TestFindDependencies_NoGap(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	entries := []struct {
		Expr    string
		Command string
	}{
		{"0 * * * *", "a.sh"},
		{"30 * * * *", "b.sh"},
	}

	deps, err := FindDependencies(entries, start, end, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %d", len(deps))
	}
}

func TestFindDependencies_InvalidGap(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	_, err := FindDependencies(nil, start, end, 0)
	if err == nil {
		t.Fatal("expected error for zero gapMinutes")
	}
}

func TestFindDependencies_InvalidExpr(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	entries := []struct {
		Expr    string
		Command string
	}{
		{"not-a-cron", "bad.sh"},
	}
	_, err := FindDependencies(entries, start, end, 5)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestDetectGap_EmptySlices(t *testing.T) {
	result := detectGap([]time.Time{}, []time.Time{}, 10)
	if result != 0 {
		t.Errorf("expected 0 for empty slices, got %d", result)
	}
}
