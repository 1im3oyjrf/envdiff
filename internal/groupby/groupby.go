package groupby

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Group represents a named collection of env key-value pairs.
type Group struct {
	Name    string
	Entries []Entry
}

// Entry is a single key-value pair within a group.
type Entry struct {
	Key   string
	Value string
}

// GroupByPrefix groups env entries from a file by their key prefix (e.g. DB_, AWS_).
// Keys without an underscore are placed in the "OTHER" group.
func GroupByPrefix(path string) ([]Group, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return groupByPrefix(f)
}

func groupByPrefix(r io.Reader) ([]Group, error) {
	groupMap := map[string][]Entry{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		prefix := "OTHER"
		if i := strings.IndexByte(key, '_'); i > 0 {
			prefix = key[:i]
		}
		groupMap[prefix] = append(groupMap[prefix], Entry{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(groupMap))
	for name := range groupMap {
		names = append(names, name)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		entries := groupMap[name]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
		groups = append(groups, Group{Name: name, Entries: entries})
	}
	return groups, nil
}
