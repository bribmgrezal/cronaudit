package audit

import (
	"fmt"
	"strings"
)

// HostSimilarityResult holds the similarity audit result for a single host.
type HostSimilarityResult struct {
	Host   string
	Report *SimilarityReport
	Err    error
}

// AuditHostsSimilarity runs similarity audits across multiple hosts/files.
func AuditHostsSimilarity(hosts map[string]string, threshold float64) []HostSimilarityResult {
	results := make([]HostSimilarityResult, 0, len(hosts))
	for host, path := range hosts {
		report, err := AuditSimilarity(path, threshold)
		results = append(results, HostSimilarityResult{
			Host:   host,
			Report: report,
			Err:    err,
		})
	}
	return results
}

// FormatSimilaritySummary returns a multi-host similarity summary string.
func FormatSimilaritySummary(results []HostSimilarityResult) string {
	var sb strings.Builder
	totalGroups := 0
	hostsWithGroups := 0

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&sb, "[%s] ERROR: %v\n", r.Host, r.Err)
			continue
		}
		g := len(r.Report.Groups)
		totalGroups += g
		if g > 0 {
			hostsWithGroups++
		}
		fmt.Fprintf(&sb, "[%s] %d similar group(s)\n", r.Host, g)
	}

	fmt.Fprintf(&sb, "---\nTotal similar groups: %d across %d host(s)",
		totalGroups, hostsWithGroups)
	return strings.TrimRight(sb.String(), "\n")
}
