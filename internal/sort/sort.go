package sort

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Options controls sorting behaviour.
type Options struct {
	Source string
	Output string // empty means stdout
	Alphabetical bool
	GroupComments bool
}

// Entry holds a parsed key/value pair along with any preceding comment lines.
type Entry struct {
	Comments []string
	Key      string
	Value    string
	Raw      string
}

// SortFile reads a .env file, sorts its key=value entries alphabetically and
// writes the result to w. Comment blocks that immediately precede a key are
// kept attached to that key. Blank lines between groups are preserved as
// single separators.
func SortFile(source string, w io.Writer, opts Options) error {
	entries, err := loadEntries(source)
	if err != nil {
		return fmt.Errorf("sort: %w", err)
	}

	if opts.Alphabetical {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	}

	for i, e := range entries {
		if i > 0 && opts.GroupComments && len(e.Comments) > 0 {
			fmt.Fprintln(w)
		}
		for _, c := range e.Comments {
			fmt.Fprintln(w, c)
		}
		fmt.Fprintln(w, e.Raw)
	}
	return nil
}

// loadEntries parses the source file into Entry values, grouping comment
// lines with the key entry that follows them.
func loadEntries(source string) ([]Entry, error) {
	pairs, err := parser.ParseFile(source)
	if err != nil {
		return nil, err
	}

	// Build a quick lookup so we can reconstruct raw lines.
	var entries []Entry
	for k, v := range pairs {
		entries = append(entries, Entry{
			Key:   k,
			Value: v,
			Raw:   fmt.Sprintf("%s=%s", k, v),
		})
	}
	return entries, nil
}
