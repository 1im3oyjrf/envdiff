package envequal

import (
	"fmt"
	"io"
	"sort"

	"github.com/envdiff/internal/parser"
)

// Result holds the outcome of comparing two env files for equality.
type Result struct {
	Equal        bool
	OnlyInA      []string
	OnlyInB      []string
	Mismatched   []string
	FileA        string
	FileB        string
}

// CheckEqual compares two .env files and returns a Result indicating
// whether they are exactly equal in keys and values.
func CheckEqual(fileA, fileB string) (Result, error) {
	envA, err := parser.ParseFile(fileA)
	if err != nil {
		return Result{}, fmt.Errorf("reading %s: %w", fileA, err)
	}
	envB, err := parser.ParseFile(fileB)
	if err != nil {
		return Result{}, fmt.Errorf("reading %s: %w", fileB, err)
	}

	r := Result{FileA: fileA, FileB: fileB}

	for k, vA := range envA {
		if vB, ok := envB[k]; !ok {
			r.OnlyInA = append(r.OnlyInA, k)
		} else if vA != vB {
			r.Mismatched = append(r.Mismatched, k)
		}
	}
	for k := range envB {
		if _, ok := envA[k]; !ok {
			r.OnlyInB = append(r.OnlyInB, k)
		}
	}

	sort.Strings(r.OnlyInA)
	sort.Strings(r.OnlyInB)
	sort.Strings(r.Mismatched)

	r.Equal = len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.Mismatched) == 0
	return r, nil
}

// WriteText writes a human-readable equality report to w.
func WriteText(w io.Writer, r Result) {
	if r.Equal {
		fmt.Fprintf(w, "✔ %s and %s are identical\n", r.FileA, r.FileB)
		return
	}
	fmt.Fprintf(w, "✘ %s and %s differ\n", r.FileA, r.FileB)
	for _, k := range r.OnlyInA {
		fmt.Fprintf(w, "  only in %s: %s\n", r.FileA, k)
	}
	for _, k := range r.OnlyInB {
		fmt.Fprintf(w, "  only in %s: %s\n", r.FileB, k)
	}
	for _, k := range r.Mismatched {
		fmt.Fprintf(w, "  value mismatch: %s\n", k)
	}
}
