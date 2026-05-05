package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/mask"
)

// WriteText writes a human-readable text report of the diff result to the given writer.
func WriteText(w io.Writer, result diff.Result, opts Options) error {
	if len(result.MissingInTarget) == 0 && len(result.MissingInSource) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(w, "✓ No differences found. Environments are in sync.")
		return nil
	}

	if len(result.MissingInTarget) > 0 {
		fmt.Fprintf(w, "Missing in target (%d):\n", len(result.MissingInTarget))
		for _, key := range result.MissingInTarget {
			fmt.Fprintf(w, "  - %s\n", key)
		}
	}

	if len(result.MissingInSource) > 0 {
		fmt.Fprintf(w, "Missing in source (%d):\n", len(result.MissingInSource))
		for _, key := range result.MissingInSource {
			fmt.Fprintf(w, "  + %s\n", key)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintf(w, "Mismatched values (%d):\n", len(result.Mismatched))
		for _, m := range result.Mismatched {
			srcVal := mask.ApplyMask(m.Key, m.SourceValue, opts.MaskSecrets)
			tgtVal := mask.ApplyMask(m.Key, m.TargetValue, opts.MaskSecrets)
			fmt.Fprintf(w, "  ~ %s\n", m.Key)
			fmt.Fprintf(w, "      source: %s\n", srcVal)
			fmt.Fprintf(w, "      target: %s\n", tgtVal)
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 40))
	fmt.Fprintln(w, buildSummary(result))
	return nil
}
