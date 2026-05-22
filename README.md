# cronaudit

> Parses and validates crontab files across hosts, reporting conflicts and missed schedules.

---

## Installation

```bash
go install github.com/yourorg/cronaudit@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/cronaudit.git && cd cronaudit && go build ./...
```

---

## Usage

```bash
# Audit a single crontab file
cronaudit --file /etc/crontab

# Audit crontab files across multiple hosts
cronaudit --hosts hosts.txt --crontab /etc/crontab

# Output a report in JSON format
cronaudit --file /etc/crontab --format json
```

**Example output:**

```
[CONFLICT] 0 * * * * /usr/bin/backup.sh  overlaps with  */5 * * * * /usr/bin/sync.sh
[MISSED]   0 2 * * *  last run: 72h ago (expected: 24h)
[OK]       30 6 * * 1 /usr/bin/weekly-report.sh
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to a local crontab file | — |
| `--hosts` | File containing list of remote hosts | — |
| `--crontab` | Remote path to crontab file | `/etc/crontab` |
| `--format` | Output format: `text`, `json` | `text` |
| `--threshold` | Missed schedule threshold (e.g. `24h`) | `24h` |

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

[MIT](LICENSE) © yourorg