package envfilter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how filtering is applied.
type Options struct {
	Prefix    string
	Suffix    string
	KeyPattern string
	Invert    bool
}

// Result holds the outcome of a filter operation.
type Result struct {
	File    string
	Matched map[string]string
	Total   int
}

// FilterFile reads a .env file and returns only the entries matching the given options.
func FilterFile(path string, opts Options) (Result, error) {
	entries, err := parser.ParseFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("envfilter: failed to parse %s: %w", path, err)
	}

	matched := make(map[string]string)
	for k, v := range entries {
		if matches(k, opts) {
			matched[k] = v
		}
	}

	return Result{
		File:    path,
		Matched: matched,
		Total:   len(entries),
	}, nil
}

func matches(key string, opts Options) bool {
	ok := true
	if opts.Prefix != "" {
		ok = ok && strings.HasPrefix(key, opts.Prefix)
	}
	if opts.Suffix != "" {
		ok = ok && strings.HasSuffix(key, opts.Suffix)
	}
	if opts.KeyPattern != "" {
		ok = ok && strings.Contains(strings.ToLower(key), strings.ToLower(opts.KeyPattern))
	}
	if opts.Invert {
		return !ok
	}
	return ok
}

// WriteText writes the filtered result as KEY=VALUE lines to w.
func WriteText(w io.Writer, r Result) {
	keys := make([]string, 0, len(r.Matched))
	for k := range r.Matched {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "# filtered from %s (%d/%d matched)\n", r.File, len(r.Matched), r.Total)
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, r.Matched[k])
	}
}

// WriteTextToStdout is a convenience wrapper.
func WriteTextToStdout(r Result) {
	WriteText(os.Stdout, r)
}
