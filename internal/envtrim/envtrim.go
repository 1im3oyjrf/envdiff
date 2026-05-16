package envtrim

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of trimming a single file.
type Result struct {
	File        string
	Trimmed     int
	Total       int
	DryRun      bool
	Lines       []TrimmedLine
}

// TrimmedLine represents a line whose value was trimmed.
type TrimmedLine struct {
	Key      string
	Original string
	Trimmed  string
}

// TrimFile reads the given .env file, strips leading/trailing whitespace
// from all values, and writes the result to output (or in-place if output == "").
// If dryRun is true, no files are written.
func TrimFile(source, output string, dryRun bool) (Result, error) {
	f, err := os.Open(source)
	if err != nil {
		return Result{}, fmt.Errorf("open %s: %w", source, err)
	}
	defer f.Close()

	var outLines []string
	var trimmed []TrimmedLine
	total := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		stripped := strings.TrimSpace(line)

		// Preserve blank lines and comments as-is.
		if stripped == "" || strings.HasPrefix(stripped, "#") {
			outLines = append(outLines, line)
			continue
		}

		eqIdx := strings.IndexByte(stripped, '=')
		if eqIdx < 0 {
			outLines = append(outLines, line)
			continue
		}

		total++
		key := stripped[:eqIdx]
		val := stripped[eqIdx+1:]
		trimmedVal := strings.TrimSpace(val)

		if trimmedVal != val {
			trimmed = append(trimmed, TrimmedLine{
				Key:      key,
				Original: val,
				Trimmed:  trimmedVal,
			})
		}
		outLines = append(outLines, key+"="+trimmedVal)
	}

	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("scan %s: %w", source, err)
	}

	result := Result{
		File:    source,
		Trimmed: len(trimmed),
		Total:   total,
		DryRun:  dryRun,
		Lines:   trimmed,
	}

	if dryRun {
		return result, nil
	}

	dest := output
	if dest == "" {
		dest = source
	}

	out := strings.Join(outLines, "\n") + "\n"
	if err := os.WriteFile(dest, []byte(out), 0644); err != nil {
		return result, fmt.Errorf("write %s: %w", dest, err)
	}

	return result, nil
}
