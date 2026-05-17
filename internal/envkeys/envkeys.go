package envkeys

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the output of a ListKeys operation.
type Result struct {
	File string
	Keys []string
}

// Options controls ListKeys behaviour.
type Options struct {
	Sorted bool
	ValuesOnly bool
	PrefixFilter string
}

// ListKeys returns all keys found in the given .env file.
func ListKeys(path string, opts Options) (Result, error) {
	env, err := parser.ParseFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("envkeys: parse %q: %w", path, err)
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		if opts.PrefixFilter != "" && len(k) < len(opts.PrefixFilter) {
			continue
		}
		if opts.PrefixFilter != "" && k[:len(opts.PrefixFilter)] != opts.PrefixFilter {
			continue
		}
		keys = append(keys, k)
	}

	if opts.Sorted {
		sort.Strings(keys)
	}

	return Result{File: path, Keys: keys}, nil
}

// WriteText writes the key list in human-readable form to w.
func WriteText(w io.Writer, r Result, opts Options) {
	env, _ := parser.ParseFile(r.File)
	fmt.Fprintf(w, "File: %s (%d keys)\n", r.File, len(r.Keys))
	for _, k := range r.Keys {
		if opts.ValuesOnly {
			fmt.Fprintf(w, "  %s = %s\n", k, env[k])
		} else {
			fmt.Fprintf(w, "  %s\n", k)
		}
	}
}

// WriteTextToStdout is a convenience wrapper around WriteText.
func WriteTextToStdout(r Result, opts Options) {
	WriteText(os.Stdout, r, opts)
}
