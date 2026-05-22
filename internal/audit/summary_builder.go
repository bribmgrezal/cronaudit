package audit

import (
	"time"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// HostInput describes a single host's crontab content for summary building.
type HostInput struct {
	Host    string
	Content string
}

// BuildGlobalSummary audits each host input for conflicts and missed schedules
// within the provided window, returning a GlobalSummary.
func BuildGlobalSummary(hosts []HostInput, from, to time.Time) *GlobalSummary {
	gs := &GlobalSummary{
		GeneratedAt: time.Now().UTC(),
	}

	for _, h := range hosts {
		hs := HostSummary{
			Host:      h.Host,
			AuditedAt: time.Now().UTC(),
		}

		entries, err := parser.Parse(h.Content)
		if err != nil {
			hs.Errors = append(hs.Errors, err.Error())
			gs.Hosts = append(gs.Hosts, hs)
			continue
		}

		hs.TotalJobs = len(entries)

		// Count conflicts
		report := Audit(h.Content)
		hs.Conflicts = len(report.Conflicts)
		for _, pe := range report.ParseErrors {
			hs.Errors = append(hs.Errors, pe)
		}

		// Count missed schedules
		for _, entry := range entries {
			sched, parseErr := schedule.Parse(entry.Schedule)
			if parseErr != nil {
				continue
			}
			missed := schedule.FindMissed(sched, from, to)
			hs.Missed += len(missed)
		}

		gs.Hosts = append(gs.Hosts, hs)
	}

	return gs
}
