package envclone

import (
	"fmt"
	"os"
	"strings"
)

// Options controls the clone behaviour.
type Options struct {
	Source      string
	Dest        string
	Prefix      string // only clone keys with this prefix
	StripPrefix bool   // remove prefix from cloned keys
	Overwrite   bool   // overwrite existing keys in dest
}

// Result describes what happened during a clone operation.
type Result struct {
	Copied    []string
	Skipped   []string
	Overwrite []string
}

// CloneFile reads key/value pairs from Source and appends/merges them into Dest.
func CloneFile(opts Options) (Result, error) {
	srcLines, err := readLines(opts.Source)
	if err != nil {
		return Result{}, fmt.Errorf("read source: %w", err)
	}

	destLines, err := readLines(opts.Dest)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("read dest: %w", err)
	}

	// Build a map of existing dest keys for fast lookup.
	existing := map[string]int{} // key -> line index
	for i, line := range destLines {
		if k, _, ok := splitLine(line); ok {
			existing[k] = i
		}
	}

	var res Result
	for _, line := range srcLines {
		k, v, ok := splitLine(line)
		if !ok {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		destKey := k
		if opts.StripPrefix && opts.Prefix != "" {
			destKey = strings.TrimPrefix(k, opts.Prefix)
		}
		newLine := destKey + "=" + v
		if idx, exists := existing[destKey]; exists {
			if !opts.Overwrite {
				res.Skipped = append(res.Skipped, destKey)
				continue
			}
			destLines[idx] = newLine
			res.Overwrite = append(res.Overwrite, destKey)
		} else {
			destLines = append(destLines, newLine)
			res.Copied = append(res.Copied, destKey)
		}
	}

	content := strings.Join(destLines, "\n")
	if len(destLines) > 0 {
		content += "\n"
	}
	if err := os.WriteFile(opts.Dest, []byte(content), 0644); err != nil {
		return Result{}, fmt.Errorf("write dest: %w", err)
	}
	return res, nil
}

func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}
	return lines, nil
}

func splitLine(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	idx := strings.IndexByte(line, '=')
	if idx < 1 {
		return "", "", false
	}
	return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
}
