package search

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/envdiff/internal/parser"
)

// Result holds a single match found during a search.
type Result struct {
	File  string
	Key   string
	Value string
	Line  int
}

// Options controls search behaviour.
type Options struct {
	// KeyPattern is a substring or exact key to search for (case-insensitive).
	KeyPattern string
	// ValuePattern is an optional substring to match against values.
	ValuePattern string
	// ExactKey requires an exact (case-sensitive) key match when true.
	ExactKey bool
}

// SearchFiles searches one or more .env files for keys/values matching opts.
func SearchFiles(files []string, opts Options) ([]Result, error) {
	if opts.KeyPattern == "" && opts.ValuePattern == "" {
		return nil, fmt.Errorf("at least one of KeyPattern or ValuePattern must be set")
	}

	var results []Result

	for _, f := range files {
		matches, err := searchFile(f, opts)
		if err != nil {
			return nil, fmt.Errorf("searching %s: %w", f, err)
		}
		results = append(results, matches...)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Key < results[j].Key
	})

	return results, nil
}

func searchFile(path string, opts Options) ([]Result, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(raw), "\n")
	env, err := parser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Build a line-number index: key -> line number (1-based).
	lineIndex := buildLineIndex(lines)

	var results []Result
	for key, value := range env {
		if !matchesKey(key, opts) {
			continue
		}
		if opts.ValuePattern != "" && !strings.Contains(strings.ToLower(value), strings.ToLower(opts.ValuePattern)) {
			continue
		}
		results = append(results, Result{
			File:  path,
			Key:   key,
			Value: value,
			Line:  lineIndex[key],
		})
	}
	return results, nil
}

func matchesKey(key string, opts Options) bool {
	if opts.KeyPattern == "" {
		return true
	}
	if opts.ExactKey {
		return key == opts.KeyPattern
	}
	return strings.Contains(strings.ToLower(key), strings.ToLower(opts.KeyPattern))
}

func buildLineIndex(lines []string) map[string]int {
	index := make(map[string]int)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			index[strings.TrimSpace(parts[0])] = i + 1
		}
	}
	return index
}
