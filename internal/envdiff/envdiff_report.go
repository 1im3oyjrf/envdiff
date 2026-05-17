package envdiff

import (
	"fmt"
	"io"
	"sort"
)

// WriteText writes a human-readable diff report to w.
// If maskSecrets is true, values for secret-like keys are masked.
func WriteText(w io.Writer, result *DiffResult, maskSecrets bool) {
	if len(result.OnlyInA) == 0 && len(result.OnlyInB) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "✔ No differences found.")
		return
	}

	if len(result.OnlyInA) > 0 {
		fmt.Fprintf(w, "\n[Only in %s]\n", result.FileA)
		keys := sortedKeys(result.OnlyInA)
		for _, k := range keys {
			v := result.OnlyInA[k]
			if maskSecrets && isSecretKey(k) {
				v = "***"
			}
			fmt.Fprintf(w, "  %s=%s\n", k, v)
		}
	}

	if len(result.OnlyInB) > 0 {
		fmt.Fprintf(w, "\n[Only in %s]\n", result.FileB)
		keys := sortedKeys(result.OnlyInB)
		for _, k := range keys {
			v := result.OnlyInB[k]
			if maskSecrets && isSecretKey(k) {
				v = "***"
			}
			fmt.Fprintf(w, "  %s=%s\n", k, v)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintln(w, "\n[Mismatched values]")
		sorted := make([]MismatchedEntry, len(result.Mismatched))
		copy(sorted, result.Mismatched)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Key < sorted[j].Key
		})
		for _, m := range sorted {
			va, vb := m.ValueA, m.ValueB
			if maskSecrets && isSecretKey(m.Key) {
				va, vb = "***", "***"
			}
			fmt.Fprintf(w, "  %s: %q (A) vs %q (B)\n", m.Key, va, vb)
		}
	}

	WriteSummary(w, result)
}

// WriteSummary prints a summary line of the diff result.
func WriteSummary(w io.Writer, result *DiffResult) {
	fmt.Fprintf(w, "\nSummary: %d only-in-A, %d only-in-B, %d mismatched\n",
		len(result.OnlyInA), len(result.OnlyInB), len(result.Mismatched))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
