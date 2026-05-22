package audit

import (
	"fmt"
	"strings"
)

// HostTagReport pairs a hostname with its TagReport.
type HostTagReport struct {
	Host   string
	Report TagReport
}

// AuditHostsTags runs AuditTags against each host's crontab content and
// returns per-host results plus an aggregated tag frequency summary.
func AuditHostsTags(hosts map[string]string) ([]HostTagReport, map[string]int) {
	freq := make(map[string]int)
	results := make([]HostTagReport, 0, len(hosts))

	for host, content := range hosts {
		r := AuditTags(host, content)
		results = append(results, HostTagReport{Host: host, Report: r})
		for _, e := range r.Entries {
			for _, tg := range e.Tags {
				freq[tg.Key]++
			}
		}
	}

	return results, freq
}

// FormatTagSummary returns a concise multi-host tag frequency table.
func FormatTagSummary(freq map[string]int) string {
	if len(freq) == 0 {
		return "No tag data available.\n"
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, "Tag frequency across all hosts:")
	for key, count := range freq {
		fmt.Fprintf(&sb, "  %-24s %d\n", key, count)
	}
	return sb.String()
}
