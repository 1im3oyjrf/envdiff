package template

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Entry represents a single key in a template file with optional comment.
type Entry struct {
	Key     string
	Comment string
}

// GenerateFromFiles produces a merged template .env file from one or more source files.
// Keys are deduplicated and sorted. Values are replaced with empty placeholders.
func GenerateFromFiles(paths []string) ([]Entry, error) {
	seen := make(map[string]string) // key -> comment

	for _, path := range paths {
		entries, err := extractEntries(path)
		if err != nil {
			return nil, fmt.Errorf("template: reading %s: %w", path, err)
		}
		for _, e := range entries {
			if _, exists := seen[e.Key]; !exists {
				seen[e.Key] = e.Comment
			}
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]Entry, 0, len(keys))
	for _, k := range keys {
		result = append(result, Entry{Key: k, Comment: seen[k]})
	}
	return result, nil
}

// Write writes a template .env file to the given path.
func Write(path string, entries []Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("template: creating file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range entries {
		if e.Comment != "" {
			fmt.Fprintf(w, "# %s\n", e.Comment)
		}
		fmt.Fprintf(w, "%s=\n", e.Key)
	}
	return w.Flush()
}

func extractEntries(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	var pendingComment string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			pendingComment = ""
			continue
		}
		if strings.HasPrefix(line, "#") {
			pendingComment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}
		if idx := strings.IndexByte(line, '='); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			entries = append(entries, Entry{Key: key, Comment: pendingComment})
		}
		pendingComment = ""
	}
	return entries, scanner.Err()
}
