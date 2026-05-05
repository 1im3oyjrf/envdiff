package report

import (
	"encoding/json"
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/mask"
)

// JSONReport is the structured output for JSON format.
type JSONReport struct {
	Summary       string            `json:"summary"`
	MissingKeys   []string          `json:"missing_keys"`
	ExtraKeys     []string          `json:"extra_keys"`
	Mismatched    []MismatchedEntry `json:"mismatched_values"`
	TotalIssues   int               `json:"total_issues"`
}

// MismatchedEntry represents a single key with differing values.
type MismatchedEntry struct {
	Key         string `json:"key"`
	SourceValue string `json:"source_value"`
	TargetValue string `json:"target_value"`
}

// WriteJSON writes the diff result as a JSON object to w.
func WriteJSON(w io.Writer, result diff.Result, opts Options) error {
	report := JSONReport{
		Summary:     buildSummary(result),
		MissingKeys: result.MissingInTarget,
		ExtraKeys:   result.MissingInSource,
		TotalIssues: len(result.MissingInTarget) + len(result.MissingInSource) + len(result.Mismatched),
	}

	if report.MissingKeys == nil {
		report.MissingKeys = []string{}
	}
	if report.ExtraKeys == nil {
		report.ExtraKeys = []string{}
	}

	for _, m := range result.Mismatched {
		report.Mismatched = append(report.Mismatched, MismatchedEntry{
			Key:         m.Key,
			SourceValue: mask.ApplyMask(m.Key, m.SourceValue, opts.MaskSecrets),
			TargetValue: mask.ApplyMask(m.Key, m.TargetValue, opts.MaskSecrets),
		})
	}
	if report.Mismatched == nil {
		report.Mismatched = []MismatchedEntry{}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
