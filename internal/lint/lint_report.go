package lint

import (
	"fmt"
	"io"
	"strings"
)

// FileResult holds lint results for a single file.
type FileResult struct {
	File   string
	Issues []Issue
}

// WriteSummary writes a human-readable lint summary to the provided writer.
func WriteSummary(w io.Writer, results []FileResult) {
	totalIssues := 0
	for _, r := range results {
		totalIssues += len(r.Issues)
	}

	if totalIssues == 0 {
		fmt.Fprintln(w, "✔ No lint issues found.")
		return
	}

	for _, r := range results {
		if len(r.Issues) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n%s\n", r.File)
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", len(r.File)))
		for _, issue := range r.Issues {
			lineInfo := ""
			if issue.Line > 0 {
				lineInfo = fmt.Sprintf("line %d: ", issue.Line)
			}
			fmt.Fprintf(w, "  [%s] %s%s\n", severityLabel(issue.Severity), lineInfo, issue.Message)
		}
	}

	fileWord := "file"
	if len(results) != 1 {
		fileWord = "files"
	}
	issueWord := "issue"
	if totalIssues != 1 {
		issueWord = "issues"
	}
	fmt.Fprintf(w, "\nFound %d %s across %d %s.\n", totalIssues, issueWord, len(results), fileWord)
}

func severityLabel(s Severity) string {
	switch s {
	case SeverityError:
		return "ERROR"
	case SeverityWarning:
		return "WARN"
	default:
		return "INFO"
	}
}
