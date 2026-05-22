package audit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestGlobalSummary_Totals(t *testing.T) {
	gs := &GlobalSummary{
		Hosts: []HostSummary{
			{Host: "host-a", Conflicts: 2, Missed: 3, TotalJobs: 5},
			{Host: "host-b", Conflicts: 0, Missed: 1, TotalJobs: 2},
			{Host: "host-c", Conflicts: 0, Missed: 0, TotalJobs: 4},
		},
	}

	if got := gs.TotalConflicts(); got != 2 {
		t.Errorf("TotalConflicts = %d, want 2", got)
	}
	if got := gs.TotalMissed(); got != 4 {
		t.Errorf("TotalMissed = %d, want 4", got)
	}
	if got := gs.HostsWithIssues(); got != 2 {
		t.Errorf("HostsWithIssues = %d, want 2", got)
	}
}

func TestGlobalSummary_Format(t *testing.T) {
	gs := &GlobalSummary{
		GeneratedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Hosts: []HostSummary{
			{Host: "web-01", Conflicts: 1, Missed: 2, TotalJobs: 6},
			{Host: "db-01", Conflicts: 0, Missed: 0, TotalJobs: 3},
		},
	}

	var buf bytes.Buffer
	gs.Format(&buf)
	out := buf.String()

	for _, want := range []string{
		"CronAudit Global Summary",
		"web-01",
		"db-01",
		"2024-01-15",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output missing %q", want)
		}
	}
}

func TestBuildGlobalSummary(t *testing.T) {
	now := time.Now().UTC()
	from := now.Add(-2 * time.Hour)
	to := now

	hosts := []HostInput{
		{
			Host:    "host-a",
			Content: "* * * * * /usr/bin/backup\n* * * * * /usr/bin/sync\n",
		},
		{
			Host:    "host-b",
			Content: "0 2 * * * /usr/bin/nightly\n",
		},
		{
			Host:    "host-err",
			Content: "not a valid crontab line!!!@@##",
		},
	}

	gs := BuildGlobalSummary(hosts, from, to)

	if gs == nil {
		t.Fatal("expected non-nil GlobalSummary")
	}
	if len(gs.Hosts) != 3 {
		t.Errorf("expected 3 host summaries, got %d", len(gs.Hosts))
	}

	// host-a has two identical schedules — should detect conflict
	var hostA HostSummary
	for _, h := range gs.Hosts {
		if h.Host == "host-a" {
			hostA = h
		}
	}
	if hostA.TotalJobs != 2 {
		t.Errorf("host-a: expected 2 jobs, got %d", hostA.TotalJobs)
	}
	if hostA.Conflicts < 1 {
		t.Errorf("host-a: expected at least 1 conflict, got %d", hostA.Conflicts)
	}
}
