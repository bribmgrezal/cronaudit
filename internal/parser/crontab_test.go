package parser

import (
	"strings"
	"testing"
)

func TestParse_ValidEntries(t *testing.T) {
	input := `
# daily backup
0 2 * * * /usr/bin/backup.sh
*/15 * * * * /usr/bin/healthcheck.sh --quiet
MAILTO=admin@example.com
`
	result := Parse(strings.NewReader(input))

	if len(result.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}

	e := result.Entries[0]
	if e.Minute != "0" || e.Hour != "2" || e.Command != "/usr/bin/backup.sh" {
		t.Errorf("unexpected first entry: %+v", e)
	}

	e2 := result.Entries[1]
	if e2.Minute != "*/15" || e2.Command != "/usr/bin/healthcheck.sh --quiet" {
		t.Errorf("unexpected second entry: %+v", e2)
	}
}

func TestParse_InvalidEntry(t *testing.T) {
	input := `0 2 * * /only-four-schedule-fields`
	result := Parse(strings.NewReader(input))

	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if len(result.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result.Entries))
	}
}

func TestParse_EmptyAndComments(t *testing.T) {
	input := `
# this is a comment

   # another comment
`
	result := Parse(strings.NewReader(input))

	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected 0 errors, got %v", result.Errors)
	}
}

func TestParse_LineNumbers(t *testing.T) {
	input := "# comment\nbad line\n0 0 * * * /cmd"
	result := Parse(strings.NewReader(input))

	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Line != 2 {
		t.Errorf("expected error on line 2, got line %d", result.Errors[0].Line)
	}
	if len(result.Entries) != 1 || result.Entries[0].LineNumber != 3 {
		t.Errorf("expected valid entry on line 3")
	}
}

func TestIsEnvAssignment(t *testing.T) {
	cases := []struct {
		line   string
		expect bool
	}{
		{"MAILTO=root", true},
		{"PATH=/usr/bin:/bin", true},
		{"0 * * * * cmd", false},
		{"= value", false},
	}
	for _, c := range cases {
		got := isEnvAssignment(c.line)
		if got != c.expect {
			t.Errorf("isEnvAssignment(%q) = %v, want %v", c.line, got, c.expect)
		}
	}
}
