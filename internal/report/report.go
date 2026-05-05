package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/mask"
)

// Options controls report formatting.
type Options struct {
	MaskSecrets bool
	SourceLabel string
	TargetLabel string
}

// DefaultOptions returns sensible defaults for report generation.
func DefaultOptions() Options {
	return Options{
		MaskSecrets: true,
		SourceLabel: "source",
		TargetLabel: "target",
	}
}

// Write renders a human-readable diff report to w.
func Write(w io.Writer, result diff.Result, opts Options) error {
	if result.Clean() {
		_, err := fmt.Fprintln(w, "✓ No differences found.")
		return err
	}

	if len(result.MissingInTarget) > 0 {
		fmt.Fprintf(w, "Missing in %s (%d):\n", opts.TargetLabel, len(result.MissingInTarget))
		for _, k := range result.MissingInTarget {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(result.MissingInSource) > 0 {
		fmt.Fprintf(w, "Missing in %s (%d):\n", opts.SourceLabel, len(result.MissingInSource))
		for _, k := range result.MissingInSource {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(result.Mismatched) > 0 {
		fmt.Fprintf(w, "Mismatched values (%d):\n", len(result.Mismatched))
		for _, m := range result.Mismatched {
			srcVal := mask.ApplyMask(m.Key, m.SourceValue, opts.MaskSecrets)
			tgtVal := mask.ApplyMask(m.Key, m.TargetValue, opts.MaskSecrets)
			fmt.Fprintf(w, "  ~ %s: %s=%s | %s=%s\n",
				m.Key,
				opts.SourceLabel, srcVal,
				opts.TargetLabel, tgtVal,
			)
		}
	}

	summary := buildSummary(result)
	_, err := fmt.Fprintf(w, "Summary: %s\n", summary)
	return err
}

func buildSummary(r diff.Result) string {
	parts := []string{}
	if n := len(r.MissingInTarget); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in target", n))
	}
	if n := len(r.MissingInSource); n > 0 {
		parts = append(parts, fmt.Sprintf("%d missing in source", n))
	}
	if n := len(r.Mismatched); n > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", n))
	}
	return strings.Join(parts, ", ")
}
