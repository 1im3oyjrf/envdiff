package redact

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/mask"
)

// Result holds the outcome of redacting a single file.
type Result struct {
	File        string
	LinesTotal  int
	LinesRedacted int
}

// RedactFile reads an .env file and writes a copy with secret values replaced
// by "***". Non-secret keys are written as-is. Comments and blank lines are
// preserved.
func RedactFile(src string, dst io.Writer) (Result, error) {
	f, err := os.Open(src)
	if err != nil {
		return Result{}, fmt.Errorf("redact: open %s: %w", src, err)
	}
	defer f.Close()

	result := Result{File: src}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		result.LinesTotal++

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			fmt.Fprintln(dst, line)
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintln(dst, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := parts[1]

		if mask.IsSecret(key) {
			result.LinesRedacted++
			fmt.Fprintf(dst, "%s=***\n", key)
		} else {
			fmt.Fprintf(dst, "%s=%s\n", key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("redact: scan %s: %w", src, err)
	}
	return result, nil
}
