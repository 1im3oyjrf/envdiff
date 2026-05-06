package lint

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// RunFile reads an env file from disk, lints it, and writes a human-readable
// report to the provided writer. Returns true if any errors were found.
func RunFile(filename string, w io.Writer) (bool, error) {
	lines, err := readLines(filename)
	if err != nil {
		return false, fmt.Errorf("lint: cannot read file %q: %w", filename, err)
	}

	result := CheckFile(filename, lines)

	if len(result.Issues) == 0 {
		fmt.Fprintf(w, "[lint] %s: OK\n", filename)
		return false, nil
	}

	fmt.Fprintf(w, "[lint] %s: %d issue(s) found\n", filename, len(result.Issues))
	for _, issue := range result.Issues {
		severityLabel := strings.ToUpper(issue.Severity)
		if issue.Key != "" {
			fmt.Fprintf(w, "  [%s] line %d (%s): %s\n", severityLabel, issue.Line, issue.Key, issue.Message)
		} else {
			fmt.Fprintf(w, "  [%s] line %d: %s\n", severityLabel, issue.Line, issue.Message)
		}
	}

	return result.HasErrors(), nil
}

// RunFiles lints multiple files and writes a combined report.
// Returns true if any file had lint errors.
func RunFiles(filenames []string, w io.Writer) (bool, error) {
	anyErrors := false
	for _, f := range filenames {
		hasErr, err := RunFile(f, w)
		if err != nil {
			return false, err
		}
		if hasErr {
			anyErrors = true
		}
	}
	return anyErrors, nil
}

// readLines reads all lines from a file into a string slice.
func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
