package envcheck

import (
	"fmt"
	"io"
)

// WriteText writes a human-readable envcheck report to w.
func WriteText(w io.Writer, results []Result, maskSecrets bool) {
	present, missing, mismatched := Summary(results)

	for _, r := range results {
		if r.Present && r.Match {
			fmt.Fprintf(w, "  OK        %s\n", r.Key)
			continue
		}
		fmt.Fprintf(w, "  %s\n", FormatMismatch(r, maskSecrets))
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "Summary: %d ok, %d missing, %d mismatched\n", present, missing, mismatched)
}

// WriteSummaryText writes only the summary line to w.
func WriteSummaryText(w io.Writer, results []Result) {
	present, missing, mismatched := Summary(results)
	fmt.Fprintf(w, "envcheck: %d ok, %d missing, %d mismatched\n", present, missing, mismatched)
}

// HasIssues returns true when any result is missing or mismatched.
func HasIssues(results []Result) bool {
	for _, r := range results {
		if !r.Present || !r.Match {
			return true
		}
	}
	return false
}
