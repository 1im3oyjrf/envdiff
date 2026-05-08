package setkey

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of a SetKey operation.
type Result struct {
	Key     string
	Value   string
	Updated bool // true if key existed and was updated, false if it was added
}

// SetKey sets or adds a key=value pair in the given .env file.
// If the key already exists, its value is replaced in-place.
// If it does not exist, it is appended at the end of the file.
func SetKey(path, key, value string, overwrite bool) (Result, error) {
	lines, err := readLines(path)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("reading file: %w", err)
	}

	updated := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || trimmed == "" {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) < 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) == key {
			if !overwrite {
				return Result{Key: key, Value: parts[1], Updated: false},
					fmt.Errorf("key %q already exists; use overwrite=true to replace", key)
			}
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	if err := writeLines(path, lines); err != nil {
		return Result{}, fmt.Errorf("writing file: %w", err)
	}

	return Result{Key: key, Value: value, Updated: updated}, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, err
		}
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

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
