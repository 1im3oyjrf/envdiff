package envcheck

import (
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of checking a single env key against the current environment.
type Result struct {
	Key      string
	Expected string
	Actual   string
	Present  bool
	Match    bool
}

// Options controls envcheck behaviour.
type Options struct {
	// MaskSecrets hides values when reporting mismatches.
	MaskSecrets bool
	// IgnoreCase performs case-insensitive value comparison.
	IgnoreCase bool
}

// CheckAgainstEnv compares the key/value pairs from a parsed env map against
// the actual running process environment. It returns one Result per key.
func CheckAgainstEnv(envMap map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(envMap))

	for key, expected := range envMap {
		actual, present := os.LookupEnv(key)

		var match bool
		if present {
			if opts.IgnoreCase {
				match = strings.EqualFold(actual, expected)
			} else {
				match = actual == expected
			}
		}

		results = append(results, Result{
			Key:      key,
			Expected: expected,
			Actual:   actual,
			Present:  present,
			Match:    match,
		})
	}

	sortResults(results)
	return results
}

// Summary returns counts of present/missing/mismatched results.
func Summary(results []Result) (present, missing, mismatched int) {
	for _, r := range results {
		if !r.Present {
			missing++
		} else if !r.Match {
			mismatched++
		} else {
			present++
		}
	}
	return
}

// FormatMismatch returns a human-readable description of a single mismatch.
func FormatMismatch(r Result, maskSecrets bool) string {
	expected := r.Expected
	actual := r.Actual
	if maskSecrets {
		expected = "***"
		actual = "***"
	}
	if !r.Present {
		return fmt.Sprintf("MISSING   %s (expected %q)", r.Key, expected)
	}
	return fmt.Sprintf("MISMATCH  %s: expected %q, got %q", r.Key, expected, actual)
}

func sortResults(results []Result) {
	// insertion sort to avoid importing sort for small slices, but use sort pkg for correctness
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Key < results[j-1].Key; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
