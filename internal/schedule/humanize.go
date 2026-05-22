package schedule

import "fmt"

// HumanizeResult holds a human-readable description of a cron expression.
type HumanizeResult struct {
	Expression string
	Description string
	Frequency   string
}

// Humanize converts a cron expression into a plain-English description.
func Humanize(expr string) (HumanizeResult, error) {
	fields, err := splitAndValidate(expr)
	if err != nil {
		return HumanizeResult{}, err
	}

	minute := fields[0]
	hour := fields[1]
	dom := fields[2]
	month := fields[3]
	dow := fields[4]

	freq := ClassifyFrequency(expr)

	var desc string
	switch {
	case minute == "*" && hour == "*" && dom == "*" && month == "*" && dow == "*":
		desc = "every minute"
	case isStep(minute, "*", 1) && hour == "*" && dom == "*" && month == "*" && dow == "*":
		step := stepValue(minute)
		desc = fmt.Sprintf("every %d minute(s)", step)
	case minute != "*" && hour == "*" && dom == "*" && month == "*" && dow == "*":
		desc = fmt.Sprintf("at minute %s of every hour", minute)
	case minute != "*" && hour != "*" && dom == "*" && month == "*" && dow == "*":
		desc = fmt.Sprintf("at %s:%s every day", padTwo(hour), padTwo(minute))
	case dom != "*" && month != "*":
		desc = fmt.Sprintf("at %s:%s on %s/%s", padTwo(hour), padTwo(minute), month, dom)
	case dow != "*":
		desc = fmt.Sprintf("at %s:%s on weekday %s", padTwo(hour), padTwo(minute), dow)
	default:
		desc = fmt.Sprintf("cron(%s)", expr)
	}

	return HumanizeResult{
		Expression:  expr,
		Description: desc,
		Frequency:   string(freq),
	}, nil
}

func splitAndValidate(expr string) ([]string, error) {
	if err := Validate(expr); err != nil {
		return nil, err
	}
	fields := splitFields(expr)
	return fields, nil
}

func splitFields(expr string) []string {
	var fields []string
	start := 0
	for i := 0; i <= len(expr); i++ {
		if i == len(expr) || expr[i] == ' ' || expr[i] == '\t' {
			if i > start {
				fields = append(fields, expr[start:i])
			}
			start = i + 1
		}
	}
	return fields
}

func isStep(field, base string, _ int) bool {
	_ = base
	for i, c := range field {
		if c == '/' && i > 0 {
			return true
		}
	}
	return false
}

func stepValue(field string) int {
	for i, c := range field {
		if c == '/' {
			v, err := parseInt(field[i+1:])
			if err == nil {
				return v
			}
		}
	}
	return 1
}

func padTwo(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}
