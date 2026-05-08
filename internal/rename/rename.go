package rename

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the outcome of a rename operation.
type Result struct {
	OldKey   string
	NewKey   string
	Renamed  bool
	Skipped  bool
	Reason   string
}

// Options controls rename behaviour.
type Options struct {
	DryRun    bool
	Overwrite bool
}

// RenameKey renames oldKey to newKey in the given .env file.
// Returns a Result describing what happened.
func RenameKey(path, oldKey, newKey string, opts Options) (Result, error) {
	if oldKey == "" || newKey == "" {
		return Result{}, fmt.Errorf("oldKey and newKey must not be empty")
	}

	env, err := parser.ParseFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("parse %s: %w", path, err)
	}

	if _, exists := env[oldKey]; !exists {
		return Result{OldKey: oldKey, NewKey: newKey, Skipped: true, Reason: "key not found"}, nil
	}

	if _, exists := env[newKey]; exists && !opts.Overwrite {
		return Result{OldKey: oldKey, NewKey: newKey, Skipped: true, Reason: "target key already exists"}, nil
	}

	if opts.DryRun {
		return Result{OldKey: oldKey, NewKey: newKey, Renamed: true}, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("read %s: %w", path, err)
	}

	updated := rewriteLines(string(raw), oldKey, newKey)

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return Result{}, fmt.Errorf("write %s: %w", path, err)
	}

	return Result{OldKey: oldKey, NewKey: newKey, Renamed: true}, nil
}

// rewriteLines replaces the key portion of matching assignment lines.
func rewriteLines(content, oldKey, newKey string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, oldKey+"=") || trimmed == oldKey {
			lines[i] = strings.Replace(line, oldKey, newKey, 1)
		}
	}
	return strings.Join(lines, "\n")
}
