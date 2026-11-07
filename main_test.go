package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCount int
	}{
		{
			name:      "three lines",
			content:   "INFO started\nWARN high\nERROR failed\n",
			wantCount: 3,
		},
		{
			name:      "one line without trailing newline",
			content:   "INFO started",
			wantCount: 1,
		},
		{
			name:      "empty file",
			content:   "",
			wantCount: 0,
		},
		{
			name:      "blank lines count",
			content:   "INFO started\n\nERROR failed",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "sample.log")

			err := os.WriteFile(path, []byte(tt.content), 0600)
			if err != nil {
				t.Fatalf("creating test file: %v", err)
			}

			got, err := countLines(path)
			if err != nil {
				t.Fatalf("countLines() returned error: %v", err)
			}

			if got != tt.wantCount {
				t.Errorf("countLines() = %d, want %d",
					got, tt.wantCount)
			}
		})
	}
}

func TestCountLinesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.log")

	_, err := countLines(path)
	if err == nil {
		t.Fatal("countLines() expected an error for a missing file")
	}
}
func TestCountLinesLongLine(t *testing.T) {
	content := strings.Repeat("A", 70*1024)
	path := filepath.Join(t.TempDir(), "long.log")
	err := os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		t.Fatalf("creating test file: %v", err)
	}
	got, err := countLines(path)
	if err != nil {
		t.Fatalf("countLines() returned error: %v", err)
	}
	if got != 1 {
		t.Errorf("countLines() = %d, want %d", got, 1)
	}
}
func TestCountLinesLargeFile(t *testing.T) {
	const numberOfLines = 100000
	content := strings.Repeat(
		"INFO application started\n",
		numberOfLines,
	)
	path := filepath.Join(t.TempDir(), "large.log")
	err := os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		t.Fatalf("creating large file: %v", err)
	}
	got, err := countLines(path)
	if err != nil {
		t.Fatalf("countLines() returned error: %v", err)
	}
	if got != numberOfLines {
		t.Errorf("countLines() = %d, want %d", got, numberOfLines)
	}
}

func BenchmarkCountLines(b *testing.B) {
	const numberOfLines = 100000
	content := strings.Repeat(
		"INFO application started\n",
		numberOfLines,
	)
	path := filepath.Join(b.TempDir(), "benchmark.log")
	err := os.WriteFile(path, []byte(content), 0600)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(content)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := countLines(path)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func TestParseLogLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    LogEntry
		wantErr bool
	}{
		{
			name:  "tab separated fields",
			input: "2026-11-01T10:30:00Z\tERROR\tauth\tFailed login attempt",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:00Z",
				Severity:  "ERROR",
				Service:   "auth",
				Message:   "Failed login attempt",
			},
		},
		{
			name:  "valid error log",
			input: "2026-11-01T10:30:00Z ERROR auth Failed login",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:00Z",
				Severity:  "ERROR",
				Service:   "auth",
				Message:   "Failed login",
			},
		},
		{
			name:  "valid info log",
			input: "2026-11-01T10:30:05Z INFO api Server started",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:05Z",
				Severity:  "INFO",
				Service:   "api",
				Message:   "Server started",
			},
		},
		{
			name:  "valid warn log",
			input: "2026-11-01T10:30:10Z WARN database Connection pool low",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:10Z",
				Severity:  "WARN",
				Service:   "database",
				Message:   "Connection pool low",
			},
		},
		{
			name:  "valid debug log",
			input: "2026-11-01T10:30:15Z DEBUG worker Processing job 123",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:15Z",
				Severity:  "DEBUG",
				Service:   "worker",
				Message:   "Processing job 123",
			},
		},
		{
			name:    "empty line",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid timestamp",
			input:   "invalid-timestamp ERROR auth Failed login",
			wantErr: true,
		},
		{
			name:    "unknown severity",
			input:   "2026-11-01T10:30:00Z TRACE auth Failed login",
			wantErr: true,
		},
		{
			name:    "missing service",
			input:   "2026-11-01T10:30:00Z ERROR",
			wantErr: true,
		},
		{
			name:    "missing message",
			input:   "2026-11-01T10:30:00Z ERROR auth",
			wantErr: true,
		},
		{
			name:    "missing timestamp",
			input:   "ERROR auth Failed login",
			wantErr: true,
		},
		{
			name:    "missing severity",
			input:   "2026-11-01T10:30:00Z",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "     ",
			wantErr: true,
		},
		{
			name:  "message containing multiple spaces",
			input: "2026-11-01T10:30:00Z ERROR auth Failed  login attempt",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:00Z",
				Severity:  "ERROR",
				Service:   "auth",
				Message:   "Failed  login attempt",
			},
		},
		{
			name:  "repeated spaces between fields",
			input: "2026-11-01T10:30:00Z    ERROR    auth    Failed login",
			want: LogEntry{
				Timestamp: "2026-11-01T10:30:00Z",
				Severity:  "ERROR",
				Service:   "auth",
				Message:   "Failed login",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLine(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("parseLogLine() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("parseLogLine() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
