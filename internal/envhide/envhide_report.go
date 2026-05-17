package envhide

import (
	"fmt"
	"io"
)

// WriteText writes a human-readable summary of the hide operation to w.
func WriteText(w io.Writer, result Result) {
	fmt.Fprintf(w, "File: %s\n", result.File)
	if result.HiddenCount == 0 {
		fmt.Fprintln(w, "No secret keys detected — nothing hidden.")
		return
	}
	fmt.Fprintf(w, "Hidden: %d secret value(s) replaced with placeholder.\n", result.HiddenCount)
}

// WriteSummary writes a compact one-line summary suitable for multi-file output.
func WriteSummary(w io.Writer, results []Result) {
	total := 0
	for _, r := range results {
		total += r.HiddenCount
	}
	if total == 0 {
		fmt.Fprintln(w, "envhide: no secrets found across all files.")
		return
	}
	fmt.Fprintf(w, "envhide: %d secret value(s) hidden across %d file(s).\n", total, len(results))
}
