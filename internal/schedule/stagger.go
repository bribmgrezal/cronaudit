package schedule

import (
	"fmt"
	"strconv"
	"strings"
)

// StaggerResult holds the original expression and a suggested staggered variant.
type StaggerResult struct {
	Original  string
	Staggered string
	OffsetMin int
	Reason    string
}

// Stagger takes a cron expression and an offset in minutes, returning a new
// expression shifted by that offset. This helps spread jobs that would
// otherwise all fire at the same minute (e.g., every host running at :00).
//
// Only the minute field is shifted. If the expression uses wildcards or steps
// in the minute field, a best-effort substitution is made.
func Stagger(expr string, offsetMin int) (StaggerResult, error) {
	if offsetMin < 0 || offsetMin >= 60 {
		return StaggerResult{}, fmt.Errorf("offset must be in range [0, 59], got %d", offsetMin)
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return StaggerResult{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	if err := Validate(expr); err != nil {
		return StaggerResult{}, fmt.Errorf("invalid expression: %w", err)
	}

	minuteField := fields[0]
	newMinute, reason, err := shiftMinuteField(minuteField, offsetMin)
	if err != nil {
		return StaggerResult{}, err
	}

	fields[0] = newMinute
	staggered := strings.Join(fields, " ")

	return StaggerResult{
		Original:  expr,
		Staggered: staggered,
		OffsetMin: offsetMin,
		Reason:    reason,
	}, nil
}

// shiftMinuteField shifts a single minute cron field by offsetMin.
func shiftMinuteField(field string, offset int) (string, string, error) {
	// Wildcard: replace with literal offset
	if field == "*" {
		return strconv.Itoa(offset), "wildcard replaced with fixed offset", nil
	}

	// Step expression like */5 or */15
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 {
			return "", "", fmt.Errorf("invalid step in minute field: %s", field)
		}
		// Build a list of minutes offset by the given amount
		var minutes []string
		for m := offset % step; m < 60; m += step {
			minutes = append(minutes, strconv.Itoa(m))
		}
		return strings.Join(minutes, ","), fmt.Sprintf("step %s shifted by %d", field, offset), nil
	}

	// Single literal minute
	if v, err := strconv.Atoi(field); err == nil {
		newMin := (v + offset) % 60
		return strconv.Itoa(newMin), fmt.Sprintf("minute %d shifted by %d", v, offset), nil
	}

	return "", "", fmt.Errorf("unsupported minute field format for staggering: %s", field)
}
