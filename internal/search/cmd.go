package search

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Flags holds parsed CLI flags for the search subcommand.
type Flags struct {
	Sources      []string
	KeyPattern   string
	ValuePattern string
	ExactKey     bool
}

// ParseFlags parses search subcommand arguments from args.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)

	key := fs.String("key", "", "Key substring to search for (case-insensitive)")
	val := fs.String("value", "", "Value substring to search for (case-insensitive)")
	exact := fs.Bool("exact", false, "Require exact key match")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	sources := fs.Args()
	if len(sources) == 0 {
		return nil, fmt.Errorf("at least one source .env file is required")
	}
	if *key == "" && *val == "" {
		return nil, fmt.Errorf("provide at least --key or --value")
	}

	return &Flags{
		Sources:      sources,
		KeyPattern:   *key,
		ValuePattern: *val,
		ExactKey:     *exact,
	}, nil
}

// Run executes the search command writing results to w.
func Run(args []string, w io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	opts := Options{
		KeyPattern:   flags.KeyPattern,
		ValuePattern: flags.ValuePattern,
		ExactKey:     flags.ExactKey,
	}

	results, err := SearchFiles(flags.Sources, opts)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(w, "No matches found.")
		return nil
	}

	currentFile := ""
	for _, r := range results {
		if r.File != currentFile {
			if currentFile != "" {
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, "==> %s\n", r.File)
			currentFile = r.File
		}
		value := r.Value
		if len(value) > 60 {
			value = value[:60] + "..."
		}
		fmt.Fprintf(w, "  [line %d] %s = %s\n", r.Line, r.Key, strings.TrimSpace(value))
	}

	fmt.Fprintf(w, "\n%d match(es) across %d file(s).\n", len(results), countFiles(results))
	return nil
}

func countFiles(results []Result) int {
	seen := make(map[string]struct{})
	for _, r := range results {
		seen[r.File] = struct{}{}
	}
	return len(seen)
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	if err := Run(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
