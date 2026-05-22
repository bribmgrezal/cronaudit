package audit

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// HostSource maps a hostname to a crontab file path.
type HostSource struct {
	Host string
	Path string
}

// MultiReport aggregates audit results across multiple hosts.
type MultiReport struct {
	Reports []*ConflictReport
}

// AuditHosts runs Audit concurrently for each HostSource.
func AuditHosts(sources []HostSource) (*MultiReport, []error) {
	type result struct {
		report *ConflictReport
		err    error
	}

	results := make([]result, len(sources))
	var wg sync.WaitGroup

	for i, src := range sources {
		wg.Add(1)
		go func(idx int, s HostSource) {
			defer wg.Done()
			f, err := os.Open(s.Path)
			if err != nil {
				results[idx].err = fmt.Errorf("host %s: open %s: %w", s.Host, s.Path, err)
				return
			}
			defer f.Close()
			rep, err := Audit(s.Host, f)
			results[idx] = result{report: rep, err: err}
		}(i, src)
	}
	wg.Wait()

	multi := &MultiReport{}
	var errs []error
	for _, res := range results {
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		multi.Reports = append(multi.Reports, res.report)
	}
	return multi, errs
}

// Format writes a combined report for all hosts to w.
func (m *MultiReport) Format(w io.Writer) {
	totalConflicts := 0
	for _, r := range m.Reports {
		r.Format(w)
		fmt.Fprintln(w)
		totalConflicts += len(r.Conflicts)
	}
	fmt.Fprintf(w, "--- Summary: %d host(s), %d total conflict(s) ---\n",
		len(m.Reports), totalConflicts)
}
