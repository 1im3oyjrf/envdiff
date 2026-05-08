package copy

import (
	"fmt"
	"os"
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Options configures the CopyKey operation.
type Options struct {
	SourceFile string
	DestFile   string
	Key        string
	NewKey     string // if empty, uses Key
	Overwrite  bool
}

// Result holds the outcome of a CopyKey operation.
type Result struct {
	Key      string
	NewKey   string
	Value    string
	Created  bool // true if key was newly added
	Updated  bool // true if key was overwritten
}

// CopyKey copies a key (optionally renaming it) from source to destination env file.
func CopyKey(opts Options) (Result, error) {
	srcEnv, err := parser.ParseFile(opts.SourceFile)
	if err != nil {
		return Result{}, fmt.Errorf("reading source: %w", err)
	}

	value, ok := srcEnv[opts.Key]
	if !ok {
		return Result{}, fmt.Errorf("key %q not found in source file", opts.Key)
	}

	destKey := opts.NewKey
	if destKey == "" {
		destKey = opts.Key
	}

	lines, err := readLines(opts.DestFile)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("reading destination: %w", err)
	}

	updated, newLines, existed := upsertKey(lines, destKey, value, opts.Overwrite)
	if existed && !opts.Overwrite {
		return Result{}, fmt.Errorf("key %q already exists in destination; use overwrite to replace", destKey)
	}

	if err := writeLines(opts.DestFile, newLines); err != nil {
		return Result{}, fmt.Errorf("writing destination: %w", err)
	}

	return Result{
		Key:     opts.Key,
		NewKey:  destKey,
		Value:   value,
		Created: !updated,
		Updated: updated,
	}, nil
}

func upsertKey(lines []string, key, value string, overwrite bool) (updated bool, out []string, existed bool) {
	prefix := key + "="
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			existed = true
			if overwrite {
				lines[i] = fmt.Sprintf("%s=%s", key, value)
				return true, lines, true
			}
			return false, lines, true
		}
	}
	return false, append(lines, fmt.Sprintf("%s=%s", key, value)), false
}

func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	raw := strings.Split(string(data), "\n")
	// trim trailing empty entry from final newline
	if len(raw) > 0 && raw[len(raw)-1] == "" {
		raw = raw[:len(raw)-1]
	}
	return raw, nil
}

func writeLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
