package envdiff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string
	Value string
}

// FileSummary holds the parsed entries of a single env file.
type FileSummary struct {
	Path    string
	Entries map[string]string
}

// DiffResult holds the result of comparing two FileSummary instances.
type DiffResult struct {
	OnlyInA    []string
	OnlyInB    []string
	Different  []MismatchedEntry
	Identical  []string
}

// MismatchedEntry holds a key whose value differs between two files.
type MismatchedEntry struct {
	Key    string
	ValueA string
	ValueB string
}

// Diff compares two FileSummary instances and returns a DiffResult.
func Diff(a, b FileSummary) DiffResult {
	result := DiffResult{}

	for key, valA := range a.Entries {
		valB, exists := b.Entries[key]
		if !exists {
			result.OnlyInA = append(result.OnlyInA, key)
		} else if valA != valB {
			result.Different = append(result.Different, MismatchedEntry{Key: key, ValueA: valA, ValueB: valB})
		} else {
			result.Identical = append(result.Identical, key)
		}
	}

	for key := range b.Entries {
		if _, exists := a.Entries[key]; !exists {
			result.OnlyInB = append(result.OnlyInB, key)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Identical)
	sort.Slice(result.Different, func(i, j int) bool {
		return result.Different[i].Key < result.Different[j].Key
	})

	return result
}

// HasDifferences returns true if the DiffResult contains any differences.
func (d DiffResult) HasDifferences() bool {
	return len(d.OnlyInA) > 0 || len(d.OnlyInB) > 0 || len(d.Different) > 0
}

// WriteText writes a human-readable diff summary to w.
func WriteText(w io.Writer, a, b FileSummary, result DiffResult, maskSecrets bool) {
	fmt.Fprintf(w, "Comparing: %s <-> %s\n", a.Path, b.Path)

	if !result.HasDifferences() {
		fmt.Fprintln(w, "✔ No differences found.")
		return
	}

	for _, key := range result.OnlyInA {
		fmt.Fprintf(w, "  - %-30s only in %s\n", key, a.Path)
	}
	for _, key := range result.OnlyInB {
		fmt.Fprintf(w, "  + %-30s only in %s\n", key, b.Path)
	}
	for _, m := range result.Different {
		valA, valB := m.ValueA, m.ValueB
		if maskSecrets && isSecretKey(m.Key) {
			valA = "***"
			valB = "***"
		}
		fmt.Fprintf(w, "  ~ %-30s %q vs %q\n", m.Key, valA, valB)
	}

	fmt.Fprintf(w, "Summary: %d only-in-A, %d only-in-B, %d mismatched, %d identical\n",
		len(result.OnlyInA), len(result.OnlyInB), len(result.Different), len(result.Identical))
}

func isSecretKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, kw := range []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIALS"} {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}
