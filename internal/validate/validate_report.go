package validate

import (
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable validation report to w.
func WriteText(w io.Writer, filename string, result *Result) {
	if result.OK() {
		fmt.Fprintf(w, "✔ %s — all rules passed\n", filename)
		return
	}

	fmt.Fprintf(w, "✘ %s — %d violation(s) found\n", filename, len(result.Violations))
	for _, v := range result.Violations {
		fmt.Fprintf(w, "  [%s] %s\n", v.Key, v.Message)
	}
}

// WriteSummaryText writes an aggregated summary for multiple files.
func WriteSummaryText(w io.Writer, results map[string]*Result) {
	total := 0
	files := make([]string, 0, len(results))
	for f := range results {
		files = append(files, f)
	}
	sortStrings(files)

	for _, f := range files {
		r := results[f]
		WriteText(w, f, r)
		total += len(r.Violations)
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	if total == 0 {
		fmt.Fprintln(w, "All validations passed.")
	} else {
		fmt.Fprintf(w, "Total violations: %d\n", total)
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
