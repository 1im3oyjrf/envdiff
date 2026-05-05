package report

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/mask"
)

// Options controls report output behavior.
type Options struct {
	MaskSecrets bool
	Format      string // "text" or "json"
}

// DefaultOptions returns sensible default report options.
func DefaultOptions() Options {
	return Options{
		MaskSecrets: false,
		Format:      "text",
	}
}

// Write dispatches to the appropriate report format writer.
func Write(w io.Writer, result diff.Result, opts Options) error {
	switch opts.Format {
	case "json":
		return WriteJSON(w, result, opts)
	default:
		return WriteText(w, result, opts)
	}
}

// buildSummary returns a one-line summary string for the diff result.
func buildSummary(result diff.Result) string {
	total := len(result.MissingInTarget) + len(result.MissingInSource) + len(result.Mismatched)
	if total == 0 {
		return "Summary: environments are in sync."
	}
	return fmt.Sprintf(
		"Summary: %d issue(s) found — %d missing in target, %d missing in source, %d mismatched.",
		total,
		len(result.MissingInTarget),
		len(result.MissingInSource),
		len(result.Mismatched),
	)
}

// maskResultValues returns a copy of result with secret values masked if enabled.
func maskResultValues(result diff.Result, enabled bool) diff.Result {
	if !enabled {
		return result
	}
	masked := diff.Result{
		MissingInTarget: result.MissingInTarget,
		MissingInSource: result.MissingInSource,
	}
	for _, m := range result.Mismatched {
		masked.Mismatched = append(masked.Mismatched, diff.Mismatch{
			Key:         m.Key,
			SourceValue: mask.ApplyMask(m.Key, m.SourceValue, true),
			TargetValue: mask.ApplyMask(m.Key, m.TargetValue, true),
		})
	}
	return masked
}
