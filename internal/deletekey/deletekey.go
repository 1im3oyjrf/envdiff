package deletekey

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Options configures the DeleteKey operation.
type Options struct {
	File    string
	Key     string
	DryRun  bool
}

// Result holds the outcome of a DeleteKey operation.
type Result struct {
	File    string
	Key     string
	Deleted bool
	DryRun  bool
}

// DeleteKey removes a key from the given .env file.
// If DryRun is true, no changes are written to disk.
func DeleteKey(opts Options) (Result, error) {
	lines, err := readLines(opts.File)
	if err != nil {
		return Result{}, fmt.Errorf("deletekey: read %q: %w", opts.File, err)
	}

	newLines, deleted := filterKey(lines, opts.Key)
	if !deleted {
		return Result{File: opts.File, Key: opts.Key, Deleted: false}, nil
	}

	if !opts.DryRun {
		if err := writeLines(opts.File, newLines); err != nil {
			return Result{}, fmt.Errorf("deletekey: write %q: %w", opts.File, err)
		}
	}

	return Result{File: opts.File, Key: opts.Key, Deleted: true, DryRun: opts.DryRun}, nil
}

// filterKey removes lines that define the given key, returning the remaining
// lines and whether the key was found.
func filterKey(lines []string, key string) ([]string, bool) {
	prefix := key + "="
	found := false
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) || trimmed == key {
			found = true
			continue
		}
		out = append(out, line)
	}
	return out, found
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
