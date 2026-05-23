package schedule

import "fmt"

// RiskLevel represents the severity of a schedule risk.
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// RiskResult holds a risk assessment for a cron expression.
type RiskResult struct {
	Expression string
	Level      RiskLevel
	Reasons    []string
}

// AssessRisk evaluates a cron expression for operational risk factors.
func AssessRisk(expr string) (RiskResult, error) {
	if err := Validate(expr); err != nil {
		return RiskResult{}, fmt.Errorf("invalid expression: %w", err)
	}

	result := RiskResult{
		Expression: expr,
		Level:      RiskLow,
	}

	if IsHighFrequency(expr) {
		result.Reasons = append(result.Reasons, "high-frequency schedule (runs every minute or near-constantly)")
		result.Level = RiskHigh
	}

	tags := TagExpression(expr)
	for _, t := range tags {
		if t == "weekday-restricted" || t == "month-restricted" {
			result.Reasons = append(result.Reasons, fmt.Sprintf("complex restriction detected: %s", t))
			if result.Level == RiskLow {
				result.Level = RiskMedium
			}
		}
	}

	lints := Lint(expr)
	if len(lints) > 0 {
		for _, l := range lints {
			result.Reasons = append(result.Reasons, fmt.Sprintf("lint warning: %s", l))
		}
		if result.Level == RiskLow {
			result.Level = RiskMedium
		}
	}

	return result, nil
}
