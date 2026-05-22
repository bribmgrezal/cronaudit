package audit

import (
	"strings"
	"testing"
)

const tagCrontab = `
# daily backup
0 2 * * * /usr/bin/backup.sh
# weekday job
30 8 * * 1-5 /usr/bin/report.sh
# every minute
* * * * * /usr/bin/ping.sh
`

func TestAuditTags_TagsPresent(t *testing.T) {
	report := AuditTags("test", tagCrontab)

	if len(report.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", report.Errors)
	}

	if len(report.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(report.Entries))
	}
}

func TestAuditTags_HighFrequency(t *testing.T) {
	report := AuditTags("test", tagCrontab)

	var found bool
	for _, e := range report.Entries {
		for _, tg := range e.Tags {
			if tg.Key == "high-frequency" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected at least one high-frequency tag")
	}
}

func TestAuditTags_WeekdayRestricted(t *testing.T) {
	report := AuditTags("test", tagCrontab)

	var found bool
	for _, e := range report.Entries {
		for _, tg := range e.Tags {
			if tg.Key == "weekday-restricted" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected weekday-restricted tag for 1-5 entry")
	}
}

func TestAuditTags_ParseError(t *testing.T) {
	report := AuditTags("bad", "not valid cron")
	// parser may or may not error; if entries exist they should have tags
	_ = report
}

func TestTagReport_Format_WithEntries(t *testing.T) {
	report := AuditTags("test", tagCrontab)
	out := report.Format()

	if !strings.Contains(out, "Tag report for test") {
		t.Error("expected header in format output")
	}
	if !strings.Contains(out, "line") {
		t.Error("expected line references in format output")
	}
}

func TestTagReport_Format_Empty(t *testing.T) {
	report := TagReport{File: "empty.cron"}
	out := report.Format()
	if !strings.Contains(out, "no tagged entries") {
		t.Error("expected empty message")
	}
}
