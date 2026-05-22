package audit

import (
	"strings"
	"testing"
)

const validCrontab = `
# system jobs
0 * * * * /usr/bin/backup
*/5 9-17 * * 1-5 /usr/bin/sync
`

const invalidCrontab = `
60 * * * * /bad/minute
0 25 * * * /bad/hour
0 0 * * * /good/job
`

func TestAuditValidation_AllValid(t *testing.T) {
	r := AuditValidation("host1", validCrontab)
	if !r.Valid() {
		t.Errorf("expected valid report, got issues: %v", r.Issues)
	}
	if r.Host != "host1" {
		t.Errorf("expected host1, got %s", r.Host)
	}
}

func TestAuditValidation_DetectsInvalid(t *testing.T) {
	r := AuditValidation("host2", invalidCrontab)
	if r.Valid() {
		t.Fatal("expected invalid report")
	}
	if len(r.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(r.Issues))
	}
}

func TestAuditValidation_ParseError(t *testing.T) {
	// a line that is neither comment, env assignment, nor valid cron entry
	r := AuditValidation("host3", "not a valid crontab line at all!!!")
	// parse errors surface as issues
	if r.Valid() {
		t.Error("expected issues for unparseable content")
	}
}

func TestValidationReport_Format_Valid(t *testing.T) {
	r := ValidationReport{Host: "myhost"}
	out := r.Format()
	if !strings.Contains(out, "valid") {
		t.Errorf("expected 'valid' in output, got: %s", out)
	}
}

func TestValidationReport_Format_WithIssues(t *testing.T) {
	r := AuditValidation("hostX", invalidCrontab)
	out := r.Format()
	if !strings.Contains(out, "hostX") {
		t.Error("expected host name in output")
	}
	if !strings.Contains(out, "invalid") {
		t.Error("expected 'invalid' in output")
	}
	if !strings.Contains(out, "line") {
		t.Error("expected line numbers in output")
	}
}
