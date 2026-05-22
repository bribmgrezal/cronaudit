package schedule

import (
	"fmt"
	"strconv"
	"strings"
)

// Field represents a parsed cron field with its allowed range.
type Field struct {
	Raw   string
	Min   int
	Max   int
	Values []int
}

// Schedule represents a fully parsed cron schedule.
type Schedule struct {
	Minute     Field
	Hour       Field
	DayOfMonth Field
	Month      Field
	DayOfWeek  Field
}

// Parse parses a cron expression (5 fields) into a Schedule.
func Parse(expr string) (*Schedule, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	fields := []struct {
		raw      string
		min, max int
	}{
		{parts[0], 0, 59},
		{parts[1], 0, 23},
		{parts[2], 1, 31},
		{parts[3], 1, 12},
		{parts[4], 0, 7},
	}

	s := &Schedule{}
	parsed := make([]Field, 5)
	for i, f := range fields {
		vals, err := expandField(f.raw, f.min, f.max)
		if err != nil {
			return nil, fmt.Errorf("field %d (%q): %w", i+1, f.raw, err)
		}
		parsed[i] = Field{Raw: f.raw, Min: f.min, Max: f.max, Values: vals}
	}
	s.Minute = parsed[0]
	s.Hour = parsed[1]
	s.DayOfMonth = parsed[2]
	s.Month = parsed[3]
	s.DayOfWeek = parsed[4]
	return s, nil
}

// expandField expands a cron field expression into a sorted list of integers.
func expandField(field string, min, max int) ([]int, error) {
	set := map[int]struct{}{}
	for _, part := range strings.Split(field, ",") {
		if err := expandPart(part, min, max, set); err != nil {
			return nil, err
		}
	}
	result := make([]int, 0, len(set))
	for v := min; v <= max; v++ {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}
	return result, nil
}

func expandPart(part string, min, max int, set map[int]struct{}) error {
	if part == "*" {
		for i := min; i <= max; i++ {
			set[i] = struct{}{}
		}
		return nil
	}
	if strings.HasPrefix(part, "*/") {
		step, err := strconv.Atoi(part[2:])
		if err != nil || step <= 0 {
			return fmt.Errorf("invalid step %q", part)
		}
		for i := min; i <= max; i += step {
			set[i] = struct{}{}
		}
		return nil
	}
	if idx := strings.Index(part, "-"); idx != -1 {
		lo, err1 := strconv.Atoi(part[:idx])
		hi, err2 := strconv.Atoi(part[idx+1:])
		if err1 != nil || err2 != nil || lo < min || hi > max || lo > hi {
			return fmt.Errorf("invalid range %q", part)
		}
		for i := lo; i <= hi; i++ {
			set[i] = struct{}{}
		}
		return nil
	}
	v, err := strconv.Atoi(part)
	if err != nil || v < min || v > max {
		return fmt.Errorf("value %q out of range [%d-%d]", part, min, max)
	}
	set[v] = struct{}{}
	return nil
}
