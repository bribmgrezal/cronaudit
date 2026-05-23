package schedule

import (
	"fmt"
	"strings"
)

// NormalizeResult holds the original and normalized cron expression.
type NormalizeResult struct {
	Original   string
	Normalized string
	Changes    []string
}

// Normalize rewrites a cron expression into a canonical form:
// - Replaces named weekdays (sun, mon, ...) with numbers
// - Replaces named months (jan, feb, ...) with numbers
// - Strips redundant leading zeros
// - Expands @aliases to their 5-field equivalents
func Normalize(expr string) (NormalizeResult, error) {
	result := NormalizeResult{Original: expr}

	if alias, ok := resolveAlias(expr); ok {
		result.Changes = append(result.Changes, fmt.Sprintf("expanded alias %q", expr))
		expr = alias
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return result, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	weekdayNames := map[string]string{
		"sun": "0", "mon": "1", "tue": "2", "wed": "3",
		"thu": "4", "fri": "5", "sat": "6",
	}
	monthNames := map[string]string{
		"jan": "1", "feb": "2", "mar": "3", "apr": "4",
		"may": "5", "jun": "6", "jul": "7", "aug": "8",
		"sep": "9", "oct": "10", "nov": "11", "dec": "12",
	}

	fields[4] = replaceNames(fields[4], weekdayNames, &result.Changes, "weekday")
	fields[3] = replaceNames(fields[3], monthNames, &result.Changes, "month")

	for i, f := range fields {
		stripped := stripLeadingZeros(f)
		if stripped != f {
			result.Changes = append(result.Changes, fmt.Sprintf("stripped leading zeros in field %d", i+1))
			fields[i] = stripped
		}
	}

	result.Normalized = strings.Join(fields, " ")
	return result, nil
}

func resolveAlias(expr string) (string, bool) {
	aliases := map[string]string{
		"@yearly":  "0 0 1 1 *",
		"@annually": "0 0 1 1 *",
		"@monthly": "0 0 1 * *",
		"@weekly":  "0 0 * * 0",
		"@daily":   "0 0 * * *",
		"@midnight": "0 0 * * *",
		"@hourly":  "0 * * * *",
	}
	v, ok := aliases[strings.ToLower(strings.TrimSpace(expr))]
	return v, ok
}

func replaceNames(field string, names map[string]string, changes *[]string, kind string) string {
	lower := strings.ToLower(field)
	for name, num := range names {
		if strings.Contains(lower, name) {
			*changes = append(*changes, fmt.Sprintf("replaced %s name %q with %s", kind, name, num))
			lower = strings.ReplaceAll(lower, name, num)
		}
	}
	return lower
}

func stripLeadingZeros(field string) string {
	if field == "*" || field == "" {
		return field
	}
	parts := strings.FieldsFunc(field, func(r rune) bool {
		return r == ',' || r == '-' || r == '/'
	})
	result := field
	for _, p := range parts {
		if len(p) > 1 && p[0] == '0' {
			stripped := strings.TrimLeft(p, "0")
			if stripped == "" {
				stripped = "0"
			}
			result = strings.ReplaceAll(result, p, stripped)
		}
	}
	return result
}
