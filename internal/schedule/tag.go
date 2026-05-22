package schedule

import "fmt"

// Tag categorizes a cron expression with a human-readable label and a
// machine-readable key used for filtering and grouping.
type Tag struct {
	Key   string
	Label string
}

// TagExpression inspects a cron expression and returns a slice of Tag values
// describing its characteristics. It returns an error if the expression is
// invalid.
func TagExpression(expr string) ([]Tag, error) {
	if err := Validate(expr); err != nil {
		return nil, fmt.Errorf("tag: invalid expression %q: %w", expr, err)
	}

	var tags []Tag

	freq := ClassifyFrequency(expr)
	tags = append(tags, Tag{Key: "frequency", Label: string(freq)})

	if IsHighFrequency(expr) {
		tags = append(tags, Tag{Key: "high-frequency", Label: "high frequency"})
	}

	fields := splitFieldsForTag(expr)
	if fields == nil {
		return tags, nil
	}

	// weekday specificity
	if fields[4] != "*" {
		tags = append(tags, Tag{Key: "weekday-restricted", Label: "weekday restricted"})
	}

	// month specificity
	if fields[3] != "*" {
		tags = append(tags, Tag{Key: "month-restricted", Label: "month restricted"})
	}

	// runs at a fixed minute (not wildcard/step)
	if fields[0] != "*" && !isStep(fields[0]) {
		tags = append(tags, Tag{Key: "fixed-minute", Label: "fixed minute"})
	}

	return tags, nil
}

func splitFieldsForTag(expr string) []string {
	fields, err := splitFields(expr)
	if err != nil || len(fields) != 5 {
		return nil
	}
	return fields
}
