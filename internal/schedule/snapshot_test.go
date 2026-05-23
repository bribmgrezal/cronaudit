package schedule

import (
	"testing"
	"time"
)

func TestSnapshot_Basic(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	jobs := [][2]string{
		{"* * * * *", "echo hello"},
		{"0 * * * *", "echo hourly"},
	}
	res, err := Snapshot(jobs, ref, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if len(res.Entries[0].NextTimes) != 3 {
		t.Errorf("expected 3 next times, got %d", len(res.Entries[0].NextTimes))
	}
}

func TestSnapshot_InvalidExpression(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	jobs := [][2]string{
		{"not-valid", "cmd"},
	}
	_, err := Snapshot(jobs, ref, 2)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestSnapshot_InvalidN(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := Snapshot([][2]string{{"* * * * *", "cmd"}}, ref, 0)
	if err == nil {
		t.Error("expected error for n=0")
	}
}

func TestSnapshot_GeneratedAt(t *testing.T) {
	ref := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	res, err := Snapshot([][2]string{{"0 12 * * *", "daily"}}, ref, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.GeneratedAt.Equal(ref) {
		t.Errorf("expected GeneratedAt=%v, got %v", ref, res.GeneratedAt)
	}
}

func TestMergeSnapshots_Union(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := ref.Add(1 * time.Minute)
	t2 := ref.Add(2 * time.Minute)
	t3 := ref.Add(3 * time.Minute)

	a := &SnapshotResult{
		GeneratedAt: ref,
		Entries: []SnapshotEntry{
			{Expression: "* * * * *", Command: "cmd", NextTimes: []time.Time{t1, t2}},
		},
	}
	b := &SnapshotResult{
		GeneratedAt: ref,
		Entries: []SnapshotEntry{
			{Expression: "* * * * *", Command: "cmd", NextTimes: []time.Time{t2, t3}},
		},
	}

	merged := MergeSnapshots(a, b)
	if len(merged.Entries) != 1 {
		t.Fatalf("expected 1 entry after merge, got %d", len(merged.Entries))
	}
	if len(merged.Entries[0].NextTimes) != 3 {
		t.Errorf("expected 3 unique times, got %d", len(merged.Entries[0].NextTimes))
	}
}

func TestMergeSnapshots_DisjointEntries(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	a := &SnapshotResult{
		GeneratedAt: ref,
		Entries: []SnapshotEntry{
			{Expression: "* * * * *", Command: "a", NextTimes: []time.Time{ref}},
		},
	}
	b := &SnapshotResult{
		GeneratedAt: ref,
		Entries: []SnapshotEntry{
			{Expression: "0 * * * *", Command: "b", NextTimes: []time.Time{ref}},
		},
	}
	merged := MergeSnapshots(a, b)
	if len(merged.Entries) != 2 {
		t.Errorf("expected 2 disjoint entries, got %d", len(merged.Entries))
	}
}
