package lint

import (
	"fmt"
	"strings"
)

// Issue represents a lint warning or error found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
	Severity string // "warn" or "error"
}

// Result holds all lint issues found in a file.
type Result struct {
	File   string
	Issues []Issue
}

// HasErrors returns true if any issue has severity "error".
func (r Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

// CheckFile lints a parsed env map against the raw lines for context.
// It checks for duplicate keys, empty values, and invalid key formats.
func CheckFile(filename string, lines []string) Result {
	result := Result{File: filename}
	seen := make(map[string]int)

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Message:  fmt.Sprintf("invalid line (missing '='): %q", trimmed),
				Severity: "error",
			})
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		value := strings.TrimSpace(trimmed[eqIdx+1:])

		if !isValidKey(key) {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("key %q contains invalid characters", key),
				Severity: "error",
			})
		}

		if prevLine, exists := seen[key]; exists {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("duplicate key %q (first seen on line %d)", key, prevLine),
				Severity: "warn",
			})
		}
		seen[key] = lineNum

		if value == "" || value == `""` || value == "''" {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("key %q has an empty value", key),
				Severity: "warn",
			})
		}
	}

	return result
}

// isValidKey returns true if the key contains only alphanumeric characters and underscores.
func isValidKey(key string) bool {
	if key == "" {
		return false
	}
	for _, ch := range key {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}
