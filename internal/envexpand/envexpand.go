// Package envexpand resolves variable references within .env files.
// It expands values like FOO=${BAR}_suffix using the keys defined in the same file
// or optionally from the current process environment as a fallback.
package envexpand

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Options controls expansion behaviour.
type Options struct {
	// FallbackToOS allows unresolved references to be looked up in os.Environ.
	FallbackToOS bool
	// FailOnMissing returns an error when a referenced variable cannot be resolved.
	FailOnMissing bool
}

// Result holds a single expanded entry.
type Result struct {
	Key      string
	Original string
	Expanded string
	Changed  bool
}

// varRef matches ${VAR} and $VAR style references.
var varRef = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// ExpandFile reads the .env file at path, expands variable references in values,
// and returns the ordered list of results.
func ExpandFile(path string, opts Options) ([]Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envexpand: open %q: %w", path, err)
	}
	defer f.Close()

	// First pass: collect all defined keys so forward references are supported.
	env := map[string]string{}
	var orderedKeys []string
	var rawValues []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		val := strings.TrimSpace(trimmed[idx+1:])
		val = stripQuotes(val)
		env[key] = val
		orderedKeys = append(orderedKeys, key)
		rawValues = append(rawValues, val)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envexpand: read %q: %w", path, err)
	}

	// Second pass: expand each value.
	results := make([]Result, 0, len(orderedKeys))
	for i, key := range orderedKeys {
		raw := rawValues[i]
		expanded, err := expand(raw, env, opts)
		if err != nil {
			return nil, fmt.Errorf("envexpand: key %q: %w", key, err)
		}
		results = append(results, Result{
			Key:      key,
			Original: raw,
			Expanded: expanded,
			Changed:  expanded != raw,
		})
	}
	return results, nil
}

// expand replaces variable references in s using the provided env map.
func expand(s string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := varRef.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		submatches := varRef.FindStringSubmatch(match)
		varName := submatches[1]
		if varName == "" {
			varName = submatches[2]
		}
		if val, ok := env[varName]; ok {
			return val
		}
		if opts.FallbackToOS {
			if val, ok := os.LookupEnv(varName); ok {
				return val
			}
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("unresolved variable %q", varName)
			return match
		}
		return ""
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// Write outputs the expanded results to w in .env format.
func Write(w *os.File, results []Result) {
	for _, r := range results {
		fmt.Fprintf(w, "%s=%s\n", r.Key, r.Expanded)
	}
}

// stripQuotes removes surrounding single or double quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
