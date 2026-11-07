package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	initialBufferSize = 64 * 1024
	maxScannerToken   = 10 * 1024 * 1024
)

type LogEntry struct {
	Timestamp string
	Severity  string
	Service   string
	Message   string
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024) // 10MB max line size
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("reading file %q: %w", path, err)
	}
	return count, nil
}
func parseLogLine(line string) (LogEntry, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return LogEntry{}, fmt.Errorf("empty log line")
	}

	// Split into fields for validation and field extraction.
	fields := strings.Fields(line)

	if len(fields) < 1 {
		return LogEntry{}, fmt.Errorf("missing timestamp")
	}

	timestamp := fields[0]

	// Validate RFC3339 timestamp.
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		return LogEntry{}, fmt.Errorf("invalid timestamp %q: %w", timestamp, err)
	}

	if len(fields) < 2 {
		return LogEntry{}, fmt.Errorf("missing severity")
	}

	severity := fields[1]

	switch severity {
	case "DEBUG", "INFO", "WARN", "ERROR":
	default:
		return LogEntry{}, fmt.Errorf("unknown severity %q", severity)
	}

	if len(fields) < 3 {
		return LogEntry{}, fmt.Errorf("missing service")
	}

	service := fields[2]

	if len(fields) < 4 {
		return LogEntry{}, fmt.Errorf("missing message")
	}

	// Locate the original message without splitting its contents.
	messageStart := 0

	for _, field := range fields[:3] {
		index := strings.Index(line[messageStart:], field)
		messageStart += index + len(field)
	}

	message := strings.TrimSpace(line[messageStart:])

	return LogEntry{
		Timestamp: timestamp,
		Severity:  severity,
		Service:   service,
		Message:   message,
	}, nil
}
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: loglens <log-file>")
		os.Exit(1)
	}
	path := os.Args[1]
	count, err := countLines(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error counting lines:", err)
		os.Exit(1)
	}
	fmt.Printf("Total Lines: %d\n", count)
}
