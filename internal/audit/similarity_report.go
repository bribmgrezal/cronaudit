package audit

import (
	"fmt"
	"strings"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/schedule"
)

// SimilarityGroup holds a set of crontab entries that fire at similar times.
type SimilarityGroup struct {
	Entries   []parser.Entry
	Score     float64
	Threshold float64
}

// SimilarityReport holds groups of similar entries found in a crontab.
type SimilarityReport struct {
	File   string
	Groups []SimilarityGroup
}

// AuditSimilarity parses a crontab file and groups entries with overlapping
// schedules above the given similarity threshold (0.0–1.0).
func AuditSimilarity(filename string, threshold float64) (*SimilarityReport, error) {
	entries, err := parser.Parse(filename)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filename, err)
	}

	exprs := make([]string, len(entries))
	for i, e := range entries {
		exprs[i] = e.Schedule
	}

	indexGroups := schedule.GroupSimilar(exprs, threshold)

	var groups []SimilarityGroup
	for _, idxs := range indexGroups {
		if len(idxs) < 2 {
			continue
		}
		var groupEntries []parser.Entry
		for _, i := range idxs {
			groupEntries = append(groupEntries, entries[i])
		}
		// compute pairwise average score for the group
		totalScore := 0.0
		count := 0
		for a := 0; a < len(idxs); a++ {
			for b := a + 1; b < len(idxs); b++ {
				s, err := schedule.SimilarityScore(exprs[idxs[a]], exprs[idxs[b]])
				if err == nil {
					totalScore += s
					count++
				}
			}
		}
		avg := 0.0
		if count > 0 {
			avg = totalScore / float64(count)
		}
		groups = append(groups, SimilarityGroup{
			Entries:   groupEntries,
			Score:     avg,
			Threshold: threshold,
		})
	}
	return &SimilarityReport{File: filename, Groups: groups}, nil
}

// Format returns a human-readable summary of the similarity report.
func (r *SimilarityReport) Format() string {
	if len(r.Groups) == 0 {
		return fmt.Sprintf("[%s] No similar schedules found.", r.File)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] %d similar group(s) found:\n", r.File, len(r.Groups))
	for i, g := range r.Groups {
		fmt.Fprintf(&sb, "  Group %d (avg similarity: %.2f):\n", i+1, g.Score)
		for _, e := range g.Entries {
			fmt.Fprintf(&sb, "    line %d: %s  %s\n", e.Line, e.Schedule, e.Command)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
