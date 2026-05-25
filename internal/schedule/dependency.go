package schedule

import (
	"fmt"
	"time"
)

// Dependency describes a potential ordering dependency between two cron entries.
type Dependency struct {
	LeaderCommand  string
	LeaderExpr     string
	FollowerCommand string
	FollowerExpr    string
	GapMinutes      int
	Note            string
}

// FindDependencies inspects pairs of schedule expressions and commands,
// returning entries where one job consistently fires shortly before another
// (within gapMinutes) — suggesting an implicit dependency.
func FindDependencies(entries []struct {
	Expr    string
	Command string
}, windowStart time.Time, windowEnd time.Time, gapMinutes int) ([]Dependency, error) {
	if gapMinutes <= 0 {
		return nil, fmt.Errorf("gapMinutes must be positive")
	}

	type firing struct {
		times []time.Time
		expr  string
		cmd   string
	}

	firings := make([]firing, 0, len(entries))
	for _, e := range entries {
		times, err := NextInWindow(e.Expr, windowStart, windowEnd)
		if err != nil {
			return nil, fmt.Errorf("invalid expression %q: %w", e.Expr, err)
		}
		firings = append(firings, firing{times: times, expr: e.Expr, cmd: e.Command})
	}

	var deps []Dependency
	for i := 0; i < len(firings); i++ {
		for j := 0; j < len(firings); j++ {
			if i == j {
				continue
			}
			gap := detectGap(firings[i].times, firings[j].times, gapMinutes)
			if gap > 0 {
				deps = append(deps, Dependency{
					LeaderCommand:   firings[i].cmd,
					LeaderExpr:      firings[i].expr,
					FollowerCommand: firings[j].cmd,
					FollowerExpr:    firings[j].expr,
					GapMinutes:      gap,
					Note:            fmt.Sprintf("leader fires ~%d min before follower", gap),
				})
			}
		}
	}
	return deps, nil
}

// detectGap returns the median gap in minutes if leader consistently fires
// before follower within gapMinutes; returns 0 if no such pattern exists.
func detectGap(leader, follower []time.Time, maxGap int) int {
	if len(leader) == 0 || len(follower) == 0 {
		return 0
	}
	var gaps []int
	for _, l := range leader {
		for _, f := range follower {
			diff := int(f.Sub(l).Minutes())
			if diff > 0 && diff <= maxGap {
				gaps = append(gaps, diff)
				break
			}
		}
	}
	if len(gaps) < 2 {
		return 0
	}
	sum := 0
	for _, g := range gaps {
		sum += g
	}
	return sum / len(gaps)
}
