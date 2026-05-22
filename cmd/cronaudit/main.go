package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/cronaudit/internal/audit"
)

func main() {
	var (
		hostsFile  = flag.String("hosts", "", "path to file listing host:crontab_path pairs (one per line)")
		crontabFile = flag.String("file", "", "path to a single crontab file to audit")
		windowStart = flag.String("from", "", "start of missed-schedule window (RFC3339, e.g. 2024-01-01T00:00:00Z)")
		windowEnd   = flag.String("to", "", "end of missed-schedule window (RFC3339)")
		showSummary = flag.Bool("summary", false, "print global summary at the end")
	)
	flag.Parse()

	if *crontabFile == "" && *hostsFile == "" {
		fmt.Fprintln(os.Stderr, "error: provide -file or -hosts")
		flag.Usage()
		os.Exit(1)
	}

	var from, to time.Time
	var checkMissed bool
	if *windowStart != "" && *windowEnd != "" {
		var err error
		from, err = time.Parse(time.RFC3339, *windowStart)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing -from: %v\n", err)
			os.Exit(1)
		}
		to, err = time.Parse(time.RFC3339, *windowEnd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing -to: %v\n", err)
			os.Exit(1)
		}
		checkMissed = true
	}

	if *crontabFile != "" {
		runSingle(*crontabFile, checkMissed, from, to, *showSummary)
		return
	}

	runMulti(*hostsFile, checkMissed, from, to, *showSummary)
}

func runSingle(path string, checkMissed bool, from, to time.Time, showSummary bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
		os.Exit(1)
	}

	report := audit.Audit(path, string(data))
	fmt.Println(report.Format())

	if checkMissed {
		mr := audit.AuditMissed(path, string(data), from, to)
		fmt.Println(mr.Format())
	}

	if showSummary {
		results := []audit.HostResult{{Host: path, Report: report}}
		summary := audit.BuildGlobalSummary(results)
		fmt.Println(summary.Format())
	}
}

func runMulti(hostsFile string, checkMissed bool, from, to time.Time, showSummary bool) {
	hosts, err := audit.LoadHostsFile(hostsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading hosts file: %v\n", err)
		os.Exit(1)
	}

	results := audit.AuditHosts(hosts)
	for _, r := range results {
		fmt.Println(r.Report.Format())
	}

	if checkMissed {
		missedResults := audit.AuditHostsMissed(hosts, from, to)
		fmt.Println(audit.SummaryMissed(missedResults))
	}

	if showSummary {
		summary := audit.BuildGlobalSummary(results)
		fmt.Println(summary.Format())
	}
}
