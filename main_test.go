package main

import (
	"os"
	"path/filepath"
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
