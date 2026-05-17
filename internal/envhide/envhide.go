package envhide

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/mask"
)

// Result holds the outcome of hiding secret values in a file.
type Result struct {
	File        string
	HiddenCount int
	Lines       []string
}

// HideSecrets reads an env file and replaces secret values with a placeholder.
// Non-secret keys and comments are preserved as-is.
func HideSecrets(path, placeholder string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("envhide: open %q: %w", path, err)
	}
	defer f.Close()

	if placeholder == "" {
		placeholder = "***"
	}

	var lines []string
	hidden := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			lines = append(lines, line)
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			lines = append(lines, line)
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if mask.IsSecret(key) {
			lines = append(lines, key+"="+placeholder)
			hidden++
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("envhide: scan %q: %w", path, err)
	}

	return Result{File: path, HiddenCount: hidden, Lines: lines}, nil
}

// Write writes the result lines to the given writer or stdout if path is empty.
func Write(result Result, outPath string) error {
	var w *os.File
	if outPath == "" {
		w = os.Stdout
	} else {
		var err error
		w, err = os.Create(outPath)
		if err != nil {
			return fmt.Errorf("envhide: create %q: %w", outPath, err)
		}
		defer w.Close()
	}
	for _, l := range result.Lines {
		fmt.Fprintln(w, l)
	}
	return nil
}
