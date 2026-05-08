package stats

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// FileStats holds statistical information about a parsed .env file.
type FileStats struct {
	FilePath      string
	TotalKeys     int
	EmptyValues   int
	SecretKeys    int
	CommentLines  int
	BlankLines    int
	UniqueValues  int
	DuplicateKeys []string
}

// Analyze parses a .env file and returns stats about its contents.
func Analyze(filePath string) (FileStats, error) {
	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return FileStats{}, fmt.Errorf("stats: failed to parse %s: %w", filePath, err)
	}

	fs := FileStats{FilePath: filePath}
	seen := make(map[string]int)
	valueSet := make(map[string]struct{})
	secretKeywords := []string{"SECRET", "PASSWORD", "PASS", "TOKEN", "KEY", "PRIVATE", "API_KEY", "AUTH"}

	for key, val := range entries {
		fs.TotalKeys++
		seen[key]++

		if val == "" {
			fs.EmptyValues++
		} else {
			valueSet[val] = struct{}{}
		}

		for _, kw := range secretKeywords {
			if containsKeyword(key, kw) {
				fs.SecretKeys++
				break
			}
		}
	}

	for k, count := range seen {
		if count > 1 {
			fs.DuplicateKeys = append(fs.DuplicateKeys, k)
		}
	}
	sort.Strings(fs.DuplicateKeys)
	fs.UniqueValues = len(valueSet)

	return fs, nil
}

func containsKeyword(key, keyword string) bool {
	upper := toUpper(key)
	return len(upper) >= len(keyword) && contains(upper, keyword)
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// WriteText writes a human-readable stats summary to w.
func WriteText(w io.Writer, fs FileStats) {
	fmt.Fprintf(w, "Stats for: %s\n", fs.FilePath)
	fmt.Fprintf(w, "  Total keys     : %d\n", fs.TotalKeys)
	fmt.Fprintf(w, "  Empty values   : %d\n", fs.EmptyValues)
	fmt.Fprintf(w, "  Secret keys    : %d\n", fs.SecretKeys)
	fmt.Fprintf(w, "  Unique values  : %d\n", fs.UniqueValues)
	if len(fs.DuplicateKeys) > 0 {
		fmt.Fprintf(w, "  Duplicate keys : %v\n", fs.DuplicateKeys)
	} else {
		fmt.Fprintf(w, "  Duplicate keys : none\n")
	}
}
