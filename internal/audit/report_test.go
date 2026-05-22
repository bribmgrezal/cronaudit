package audit

import (
	"bytes"
	"strings"
	"testing"
)

const sampleCrontab = `
# Daily backup
0 2 * * * /usr/bin/backup.sh
# Also at 2am — conflicts with above
0 2 * * * /usr/bin/cleanup.sh
# Every hour
0 * * * * /usr/bin/healthcheck.sh
`

const invalidCrontab = `
bad-line-no-schedule
* * * * * /usr/bin/valid.sh
`

func TestAudit_DetectsConflicts(t *testing.T) {
	r := strings.NewReader(sampleCrontab)
	report, err := Audit("host1", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.Host != "host1" {
		t.Errorf("expected host 'host1', got %q", report.Host)
	}
	if len(report.Conflicts) == 0 {
		t.Error("expected at least one conflict, got none")
	}
}

func TestAudit_NoConflicts(t *testing.T) {
	input := "0 1 * * * /bin/job1.sh\n0 2 * * * /bin/job2.sh\n"
	r := strings.NewReader(input)
	report, err := Audit("host2", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(report.Conflicts))
	}
}

func TestAudit_ParseError(t *testing.T) {
	r := strings.NewReader(invalidCrontab)
	report, err := Audit("host3", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Errors) == 0 {
		t.Error("expected parse errors, got none")
	}
}

func TestReport_Format(t *testing.T) {
	r := strings.NewReader(sampleCrontab)
	report, err := Audit("host4", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	report.Format(&buf)
	out := buf.String()
	if !strings.Contains(out, "host4") {
		t.Error("expected host name in output")
	}
	if !strings.Contains(out, "CONFLICT") && len(report.Conflicts) > 0 {
		t.Error("expected CONFLICT in output")
	}
}

func TestTruncate(t *testing.T) {
	long := strings.Repeat("a", 60)
	out := truncate(long, 40)
	if len(out) > 40 {
		t.Errorf("expected truncated string <= 40 chars, got %d", len(out))
	}
	if !strings.HasSuffix(out, "...") {
		t.Error("expected ellipsis suffix")
	}
	short := "hello"
	if truncate(short, 40) != short {
		t.Error("short string should not be truncated")
	}
}
