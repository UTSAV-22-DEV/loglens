package main

import (
	"bufio"
	"fmt"
	"os"
)

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
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: loglens <log-file>")
		os.Exit(2)
	}
	path := os.Args[1]
	count, err := countLines(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error counting lines:", err)
		os.Exit(1)
	}
	fmt.Println("LogLens - Log Analysis CLI")
	fmt.Println("File:", path)
	fmt.Println("Lines:", count)
}
