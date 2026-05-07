package audit

import (
	"fmt"
	"io"
)

// WriteText writes a human-readable audit report to w.
func WriteText(w io.Writer, result Result) {
	if len(result.Findings) == 0 {
		fmt.Fprintln(w, "✔ Audit passed — no issues found.")
		return
	}

	fmt.Fprintln(w, "Audit Findings:")
	fmt.Fprintln(w, "---------------")
	for _, f := range result.Findings {
		icon := iconFor(f.Severity)
		fmt.Fprintf(w, "  %s [%s] %s: %s\n", icon, f.Severity, f.Key, f.Message)
	}
	fmt.Fprintln(w)
	WriteSummary(w, result)
}

// WriteSummary writes a compact summary line to w.
func WriteSummary(w io.Writer, result Result) {
	fmt.Fprintf(w, "Summary: %d finding(s) — %d error(s), %d warning(s), %d info(s)\n",
		result.Total, result.Errors, result.Warnings, result.Infos)
}

func iconFor(s Severity) string {
	switch s {
	case SeverityError:
		return "✖"
	case SeverityWarning:
		return "⚠"
	default:
		return "ℹ"
	}
}
