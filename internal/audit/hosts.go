package audit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HostEntry represents a host and the path to its crontab file.
type HostEntry struct {
	Host     string
	FilePath string
}

// HostResult pairs a HostEntry with its audit Report.
type HostResult struct {
	Host   string
	Report Report
}

// LoadHostsFile reads a file where each non-blank, non-comment line has the
// format "hostname:/path/to/crontab" and returns the parsed entries.
func LoadHostsFile(path string) ([]HostEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open hosts file: %w", err)
	}
	defer f.Close()

	var entries []HostEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("hosts file line %d: expected host:path, got %q", lineNum, line)
		}
		host := strings.TrimSpace(parts[0])
		filePath := strings.TrimSpace(parts[1])
		if host == "" || filePath == "" {
			return nil, fmt.Errorf("hosts file line %d: host and path must be non-empty", lineNum)
		}
		entries = append(entries, HostEntry{Host: host, FilePath: filePath})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning hosts file: %w", err)
	}
	return entries, nil
}
