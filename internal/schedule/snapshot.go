package schedule

import (
	"fmt"
	"sort"
	"time"
)

// SnapshotEntry holds a single scheduled job and its next N fire times.
type SnapshotEntry struct {
	Expression string
	Command    string
	NextTimes  []time.Time
}

// SnapshotResult holds all snapshot entries for a crontab.
type SnapshotResult struct {
	GeneratedAt time.Time
	Entries     []SnapshotEntry
}

// Snapshot generates the next n fire times for each (expression, command) pair
// starting from the given reference time.
func Snapshot(jobs [][2]string, ref time.Time, n int) (*SnapshotResult, error) {
	if n <= 0 {
		return nil, fmt.Errorf("snapshot: n must be positive, got %d", n)
	}

	result := &SnapshotResult{
		GeneratedAt: ref,
	}

	for _, job := range jobs {
		expr, cmd := job[0], job[1]
		times, err := NextN(expr, ref, n)
		if err != nil {
			return nil, fmt.Errorf("snapshot: invalid expression %q: %w", expr, err)
		}
		result.Entries = append(result.Entries, SnapshotEntry{
			Expression: expr,
			Command:    cmd,
			NextTimes:  times,
		})
	}

	return result, nil
}

// MergeSnapshots merges two SnapshotResults, deduplicating by expression+command,
// keeping the union of next times sorted.
func MergeSnapshots(a, b *SnapshotResult) *SnapshotResult {
	type key struct{ expr, cmd string }
	index := make(map[key]*SnapshotEntry)

	for i := range a.Entries {
		e := &a.Entries[i]
		k := key{e.Expression, e.Command}
		index[k] = &SnapshotEntry{
			Expression: e.Expression,
			Command:    e.Command,
			NextTimes:  append([]time.Time{}, e.NextTimes...),
		}
	}

	for _, e := range b.Entries {
		k := key{e.Expression, e.Command}
		if existing, ok := index[k]; ok {
			existing.NextTimes = mergeTimes(existing.NextTimes, e.NextTimes)
		} else {
			index[k] = &SnapshotEntry{
				Expression: e.Expression,
				Command:    e.Command,
				NextTimes:  append([]time.Time{}, e.NextTimes...),
			}
		}
	}

	result := &SnapshotResult{GeneratedAt: a.GeneratedAt}
	for _, v := range index {
		result.Entries = append(result.Entries, *v)
	}
	return result
}

func mergeTimes(a, b []time.Time) []time.Time {
	seen := make(map[time.Time]struct{})
	for _, t := range a {
		seen[t] = struct{}{}
	}
	for _, t := range b {
		seen[t] = struct{}{}
	}
	result := make([]time.Time, 0, len(seen))
	for t := range seen {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Before(result[j]) })
	return result
}
