package schedule

import "sort"

// SimilarityScore returns a value between 0.0 and 1.0 indicating how similar
// two cron expressions are in terms of their firing times within a 24-hour window.
func SimilarityScore(exprA, exprB string) (float64, error) {
	schedA, err := Parse(exprA)
	if err != nil {
		return 0, err
	}
	schedB, err := Parse(exprB)
	if err != nil {
		return 0, err
	}

	timesA := dailyMinutes(schedA)
	timesB := dailyMinutes(schedB)

	if len(timesA) == 0 && len(timesB) == 0 {
		return 1.0, nil
	}
	if len(timesA) == 0 || len(timesB) == 0 {
		return 0.0, nil
	}

	setA := toMinuteSet(timesA)
	setB := toMinuteSet(timesB)

	intersection := 0
	for m := range setA {
		if setB[m] {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 1.0, nil
	}
	return float64(intersection) / float64(union), nil
}

// GroupSimilar clusters cron expressions by similarity threshold.
// Returns groups of expression indices from the input slice.
func GroupSimilar(exprs []string, threshold float64) [][]int {
	n := len(exprs)
	visited := make([]bool, n)
	var groups [][]int

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		group := []int{i}
		visited[i] = true
		for j := i + 1; j < n; j++ {
			if visited[j] {
				continue
			}
			score, err := SimilarityScore(exprs[i], exprs[j])
			if err == nil && score >= threshold {
				group = append(group, j)
				visited[j] = true
			}
		}
		groups = append(groups, group)
	}
	return groups
}

// dailyMinutes returns all minute-of-day values (0..1439) a schedule fires
// across a single representative day (using day=1, month=1, weekday=1).
func dailyMinutes(s Schedule) []int {
	seen := map[int]bool{}
	for h := range s.Hours {
		for m := range s.Minutes {
			min := h*60 + m
			seen[min] = true
		}
	}
	result := make([]int, 0, len(seen))
	for m := range seen {
		result = append(result, m)
	}
	sort.Ints(result)
	return result
}

func toMinuteSet(mins []int) map[int]bool {
	s := make(map[int]bool, len(mins))
	for _, m := range mins {
		s[m] = true
	}
	return s
}
