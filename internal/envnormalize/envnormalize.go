package envnormalize

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of normalizing a single line.
type Result struct {
	Line     int
	Original string
	Normalized string
	Changed  bool
}

// NormalizeOptions controls normalization behavior.
type NormalizeOptions struct {
	TrimWhitespace bool
	LowercaseKeys  bool
	RemoveExport   bool
	QuoteValues    bool
}

// DefaultOptions returns sensible normalization defaults.
func DefaultOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimWhitespace: true,
		LowercaseKeys:  false,
		RemoveExport:   true,
		QuoteValues:    false,
	}
}

// NormalizeFile reads a .env file and returns normalized lines with change metadata.
func NormalizeFile(path string, opts NormalizeOptions) ([]Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var results []Result
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		original := scanner.Text()
		normalized := normalizeLine(original, opts)
		results = append(results, Result{
			Line:       lineNum,
			Original:   original,
			Normalized: normalized,
			Changed:    original != normalized,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	return results, nil
}

func normalizeLine(line string, opts NormalizeOptions) string {
	// Preserve blank lines and comments as-is
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		if opts.TrimWhitespace {
			return trimmed
		}
		return line
	}

	if opts.TrimWhitespace {
		line = trimmed
	}

	if opts.RemoveExport {
		line = strings.TrimPrefix(line, "export ")
	}

	eqIdx := strings.IndexByte(line, '=')
	if eqIdx < 0 {
		return line
	}

	key := line[:eqIdx]
	val := line[eqIdx+1:]

	if opts.TrimWhitespace {
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
	}

	if opts.LowercaseKeys {
		key = strings.ToLower(key)
	}

	if opts.QuoteValues && !strings.HasPrefix(val, `"`) && !strings.HasPrefix(val, `'`) && val != "" {
		val = `"` + val + `"`
	}

	return key + "=" + val
}

// Write applies normalized results to an output writer or file.
func Write(path string, results []Result) error {
	lines := make([]string, len(results))
	for i, r := range results {
		lines[i] = r.Normalized
	}
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0644)
}
