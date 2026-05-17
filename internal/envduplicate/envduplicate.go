package envduplicate

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Finding represents a duplicate value found across keys.
type Finding struct {
	Value string
	Keys  []string
}

// Result holds the outcome of a duplicate-value scan on a single file.
type Result struct {
	File     string
	Findings []Finding
}

// FindDuplicateValues scans the given .env files and returns, per file,
// any keys that share an identical non-empty value.
func FindDuplicateValues(files []string) ([]Result, error) {
	var results []Result

	for _, f := range files {
		env, err := parser.ParseFile(f)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}

		findings := detectDuplicates(env)
		results = append(results, Result{
			File:     f,
			Findings: findings,
		})
	}

	return results, nil
}

// detectDuplicates groups keys by value and returns those with more than one key.
func detectDuplicates(env map[string]string) []Finding {
	valueToKeys := make(map[string][]string)

	for k, v := range env {
		if v == "" {
			continue
		}
		valueToKeys[v] = append(valueToKeys[v], k)
	}

	var findings []Finding
	for v, keys := range valueToKeys {
		if len(keys) > 1 {
			sort.Strings(keys)
			findings = append(findings, Finding{Value: v, Keys: keys})
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Value < findings[j].Value
	})

	return findings
}

// WriteText writes a human-readable duplicate-value report to w.
func WriteText(w io.Writer, results []Result) {
	for _, r := range results {
		fmt.Fprintf(w, "File: %s\n", r.File)
		if len(r.Findings) == 0 {
			fmt.Fprintln(w, "  No duplicate values found.")
		} else {
			for _, f := range r.Findings {
				fmt.Fprintf(w, "  Value %q shared by: %v\n", f.Value, f.Keys)
			}
		}
		fmt.Fprintln(w)
	}
}

// WriteSummary writes a one-line summary per file.
func WriteSummary(w io.Writer, results []Result) {
	for _, r := range results {
		fmt.Fprintf(w, "%s: %d duplicate-value group(s)\n", r.File, len(r.Findings))
	}
}
