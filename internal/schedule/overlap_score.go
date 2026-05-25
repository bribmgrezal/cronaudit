package schedule

import "sort"

// OverlapScore represents the degree of overlap between two cron expressions.
type OverlapScore struct {
	ExprA        string
	ExprB        string
	SharedCount  int
	TotalUnique  int
	ScorePercent float64
}

// ScoreOverlap computes a percentage-based overlap score between two cron
// expressions over a 24-hour window (1440 minutes). Returns an error if
// either expression is invalid.
func ScoreOverlap(exprA, exprB string) (OverlapScore, error) {
	schedA, err := Parse(exprA)
	if err != nil {
		return OverlapScore{}, err
	}
	schedB, err := Parse(exprB)
	if err != nil {
		return OverlapScore{}, err
	}

	setA := minuteSet(schedA)
	setB := minuteSet(schedB)

	shared := 0
	for m := range setA {
		if setB[m] {
			shared++
		}
	}

	union := len(setA) + len(setB) - shared
	var pct float64
	if union > 0 {
		pct = float64(shared) / float64(union) * 100.0
	}

	return OverlapScore{
		ExprA:        exprA,
		ExprB:        exprB,
		SharedCount:  shared,
		TotalUnique:  union,
		ScorePercent: pct,
	}, nil
}

// RankOverlaps scores all pairs of expressions and returns them sorted by
// descending overlap percentage.
func RankOverlaps(exprs []string) ([]OverlapScore, error) {
	var results []OverlapScore
	for i := 0; i < len(exprs); i++ {
		for j := i + 1; j < len(exprs); j++ {
			score, err := ScoreOverlap(exprs[i], exprs[j])
			if err != nil {
				return nil, err
			}
			results = append(results, score)
		}
	}
	sort.Slice(results, func(a, b int) bool {
		return results[a].ScorePercent > results[b].ScorePercent
	})
	return results, nil
}

// minuteSet builds a set of minutes-of-day (0..1439) that a schedule fires on.
func minuteSet(s Schedule) map[int]bool {
	set := make(map[int]bool)
	for h := range s.Hours {
		for m := range s.Minutes {
			set[h*60+m] = true
		}
	}
	return set
}
