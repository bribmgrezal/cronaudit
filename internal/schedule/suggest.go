package schedule

import "fmt"

// Suggestion represents a recommended alternative cron expression.
type Suggestion struct {
	Original   string
	Suggested  string
	Reason     string
}

// Suggest analyzes a cron expression and returns improvement suggestions.
// It looks for common anti-patterns such as high-frequency schedules,
// overly broad wildcards, and non-standard step values.
func Suggest(expr string) ([]Suggestion, error) {
	if err := Validate(expr); err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	var suggestions []Suggestion

	if IsHighFrequency(expr) {
		suggestions = append(suggestions, Suggestion{
			Original:  expr,
			Suggested: suggestReducedFrequency(expr),
			Reason:    "high-frequency schedule may cause system load; consider reducing",
		})
	}

	fields, err := splitAndValidate(expr)
	if err != nil {
		return suggestions, nil
	}

	// Suggest named weekday aliases if numeric weekdays are used
	if fields[4] != "*" {
		if s, ok := suggestWeekdayAlias(expr, fields[4]); ok {
			suggestions = append(suggestions, s)
		}
	}

	return suggestions, nil
}

// suggestReducedFrequency returns a less frequent alternative.
func suggestReducedFrequency(expr string) string {
	fields, err := splitAndValidate(expr)
	if err != nil {
		return expr
	}
	// If running every minute, suggest every 5 minutes
	if fields[0] == "*" {
		fields[0] = "*/5"
		return fields[0] + " " + fields[1] + " " + fields[2] + " " + fields[3] + " " + fields[4]
	}
	return expr
}

var weekdayNames = map[string]string{
	"0": "Sun", "1": "Mon", "2": "Tue", "3": "Wed",
	"4": "Thu", "5": "Fri", "6": "Sat",
}

// suggestWeekdayAlias suggests a named alias for a single numeric weekday.
func suggestWeekdayAlias(expr, weekdayField string) (Suggestion, bool) {
	name, ok := weekdayNames[weekdayField]
	if !ok {
		return Suggestion{}, false
	}
	fields, err := splitAndValidate(expr)
	if err != nil {
		return Suggestion{}, false
	}
	suggested := fields[0] + " " + fields[1] + " " + fields[2] + " " + fields[3] + " " + name
	return Suggestion{
		Original:  expr,
		Suggested: suggested,
		Reason:    fmt.Sprintf("use named weekday %q instead of %q for readability", name, weekdayField),
	}, true
}
