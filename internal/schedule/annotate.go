package schedule

import "fmt"

// Annotation holds a human-readable label and optional metadata for a cron expression.
type Annotation struct {
	Expression  string
	Human       string
	Frequency   string
	Risk        string
	Tags        []string
	Suggestions []string
	Warnings    []string
}

// Annotate produces a full Annotation for the given cron expression by
// combining Humanize, ClassifyFrequency, AssessRisk, TagExpression, Suggest,
// and Lint into a single result.
func Annotate(expr string) (Annotation, error) {
	if err := Validate(expr); err != nil {
		return Annotation{}, fmt.Errorf("invalid expression %q: %w", expr, err)
	}

	human := Humanize(expr)

	freq := ClassifyFrequency(expr)

	risk := AssessRisk(expr)

	tag := TagExpression(expr)
	tags := []string{string(tag)}

	suggestions := Suggest(expr)

	lintWarnings := Lint(expr)
	warnings := make([]string, len(lintWarnings))
	for i, w := range lintWarnings {
		warnings[i] = w.Message
	}

	return Annotation{
		Expression:  expr,
		Human:       human,
		Frequency:   string(freq),
		Risk:        risk.Level,
		Tags:        tags,
		Suggestions: suggestions,
		Warnings:    warnings,
	}, nil
}
