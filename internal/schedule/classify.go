package schedule

import "fmt"

// Category represents a broad classification of a cron schedule.
type Category string

const (
	CategoryEveryMinute Category = "every-minute"
	CategoryHighFreq    Category = "high-frequency"
	CategoryHourly      Category = "hourly"
	CategoryDaily       Category = "daily"
	CategoryWeekly      Category = "weekly"
	CategoryMonthly     Category = "monthly"
	CategoryCustom      Category = "custom"
)

// ClassificationResult holds the category and a human-readable reason.
type ClassificationResult struct {
	Category Category
	Reason   string
}

// Classify returns a Category for the given cron expression.
// It returns an error if the expression is invalid.
func Classify(expr string) (ClassificationResult, error) {
	fields, err := splitAndValidate(expr)
	if err != nil {
		return ClassificationResult{}, fmt.Errorf("classify: %w", err)
	}

	minute := fields[0]
	hour := fields[1]
	dom := fields[2]
	month := fields[3]
	dow := fields[4]

	// Every minute
	if minute == "*" && hour == "*" && dom == "*" && month == "*" && dow == "*" {
		return ClassificationResult{Category: CategoryEveryMinute, Reason: "runs every minute"}, nil
	}

	// High frequency: step on minute with small interval
	if isStep(minute) && hour == "*" && dom == "*" && month == "*" && dow == "*" {
		step := stepValue(minute)
		if step > 0 && step <= 15 {
			return ClassificationResult{Category: CategoryHighFreq, Reason: fmt.Sprintf("runs every %d minutes", step)}, nil
		}
	}

	// Monthly: specific day of month, wildcard dow
	if dom != "*" && month == "*" && dow == "*" {
		return ClassificationResult{Category: CategoryMonthly, Reason: "runs on a specific day each month"}, nil
	}

	// Weekly: specific day of week
	if dow != "*" && dom == "*" {
		return ClassificationResult{Category: CategoryWeekly, Reason: "runs on specific weekday(s)"}, nil
	}

	// Daily: wildcard dom/dow, specific hour
	if dom == "*" && dow == "*" && month == "*" && hour != "*" && !isStep(hour) {
		return ClassificationResult{Category: CategoryDaily, Reason: "runs once or more per day at fixed hours"}, nil
	}

	// Hourly: specific minute, wildcard hour
	if minute != "*" && !isStep(minute) && hour == "*" && dom == "*" && dow == "*" {
		return ClassificationResult{Category: CategoryHourly, Reason: "runs once per hour at a fixed minute"}, nil
	}

	return ClassificationResult{Category: CategoryCustom, Reason: "custom schedule"}, nil
}
