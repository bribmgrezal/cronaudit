package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// CronEntry represents a single parsed crontab entry.
type CronEntry struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Command    string
	Raw        string
	LineNumber int
}

// ParseResult holds all entries and any parse errors from a crontab source.
type ParseResult struct {
	Entries []CronEntry
	Errors  []ParseError
}

// ParseError describes a line-level parse failure.
type ParseError struct {
	Line    int
	Content string
	Reason  string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %s (%s)", e.Line, e.Reason, e.Content)
}

// Parse reads a crontab from r and returns a ParseResult.
func Parse(r io.Reader) ParseResult {
	var result ParseResult
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Skip environment variable assignments (e.g. MAILTO=root)
		if isEnvAssignment(line) {
			continue
		}

		entry, err := parseLine(line, lineNum)
		if err != nil {
			result.Errors = append(result.Errors, *err)
			continue
		}
		result.Entries = append(result.Entries, entry)
	}

	return result
}

func isEnvAssignment(line string) bool {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return false
	}
	key := strings.TrimSpace(parts[0])
	return !strings.ContainsAny(key, " \t") && key != ""
}

func parseLine(line string, lineNum int) (CronEntry, *ParseError) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return CronEntry{}, &ParseError{
			Line:    lineNum,
			Content: line,
			Reason:  "expected at least 6 fields (5 schedule + command)",
		}
	}

	return CronEntry{
		Minute:     fields[0],
		Hour:       fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Command:    strings.Join(fields[5:], " "),
		Raw:        line,
		LineNumber: lineNum,
	}, nil
}
