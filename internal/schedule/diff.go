package schedule

import "time"

// DiffResult holds the comparison between two schedules over a window.
type DiffResult struct {
	OnlyInA []time.Time
	OnlyInB []time.Time
	Common  []time.Time
}

// DiffSchedules compares two cron expressions over a time window and returns
// times that are unique to each schedule and times they share.
func DiffSchedules(exprA, exprB string, from, to time.Time) (*DiffResult, error) {
	schedA, err := Parse(exprA)
	if err != nil {
		return nil, err
	}
	schedB, err := Parse(exprB)
	if err != nil {
		return nil, err
	}

	timesA := collectTimes(schedA, from, to)
	timesB := collectTimes(schedB, from, to)

	setA := toTimeSet(timesA)
	setB := toTimeSet(timesB)

	result := &DiffResult{}
	for _, t := range timesA {
		key := t.Unix()
		if _, ok := setB[key]; ok {
			result.Common = append(result.Common, t)
		} else {
			result.OnlyInA = append(result.OnlyInA, t)
		}
	}
	for _, t := range timesB {
		key := t.Unix()
		if _, ok := setA[key]; !ok {
			result.OnlyInB = append(result.OnlyInB, t)
		}
	}
	return result, nil
}

func collectTimes(sched [][]int, from, to time.Time) []time.Time {
	var times []time.Time
	cur := from.Truncate(time.Minute)
	for !cur.After(to) {
		if Matches(sched, cur) {
			times = append(times, cur)
		}
		cur = cur.Add(time.Minute)
	}
	return times
}

func toTimeSet(times []time.Time) map[int64]struct{} {
	m := make(map[int64]struct{}, len(times))
	for _, t := range times {
		m[t.Unix()] = struct{}{}
	}
	return m
}
