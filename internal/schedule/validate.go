package schedule

import (
	"fmt"
	"strings"
)

// ValidationError describes a problem found in a cron field.
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid %s field %q: %s", e.Field, e.Value, e.Message)
}

// fieldLimits defines the allowed [min, max] for each cron position.
var fieldLimits = []struct {
	name string
	min  int
	max  int
}{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"day-of-month", 1, 31},
	{"month", 1, 12},
	{"day-of-week", 0, 7},
}

// Validate checks a raw cron expression (5 fields) and returns all
// validation errors found, or nil if the expression is valid.
func Validate(expr string) []error {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return []error{fmt.Errorf("expected 5 fields, got %d", len(fields))}
	}

	var errs []error
	for i, f := range fields {
		lim := fieldLimits[i]
		if err := validateField(f, lim.name, lim.min, lim.max); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func validateField(field, name string, min, max int) error {
	// wildcard is always valid
	if field == "*" {
		return nil
	}
	// step on wildcard: */n
	if strings.HasPrefix(field, "*/") {
		step, err := parseInt(field[2:])
		if err != nil || step < 1 {
			return ValidationError{Field: name, Value: field, Message: "step must be a positive integer"}
		}
		return nil
	}
	// list
	for _, part := range strings.Split(field, ",") {
		if err := validatePart(part, name, min, max); err != nil {
			return err
		}
	}
	return nil
}

func validatePart(part, name string, min, max int) error {
	// range with optional step: a-b or a-b/n
	base := part
	if idx := strings.Index(part, "/"); idx >= 0 {
		step, err := parseInt(part[idx+1:])
		if err != nil || step < 1 {
			return ValidationError{Field: name, Value: part, Message: "step must be a positive integer"}
		}
		base = part[:idx]
	}
	if strings.Contains(base, "-") {
		sides := strings.SplitN(base, "-", 2)
		lo, err1 := parseInt(sides[0])
		hi, err2 := parseInt(sides[1])
		if err1 != nil || err2 != nil {
			return ValidationError{Field: name, Value: part, Message: "range bounds must be integers"}
		}
		if lo < min || hi > max || lo > hi {
			return ValidationError{Field: name, Value: part,
				Message: fmt.Sprintf("range %d-%d out of bounds [%d,%d]", lo, hi, min, max)}
		}
		return nil
	}
	v, err := parseInt(base)
	if err != nil {
		return ValidationError{Field: name, Value: part, Message: "not a valid integer"}
	}
	if v < min || v > max {
		return ValidationError{Field: name, Value: part,
			Message: fmt.Sprintf("value %d out of bounds [%d,%d]", v, min, max)}
	}
	return nil
}
