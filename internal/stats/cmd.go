package stats

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the stats command.
type Flags struct {
	Sources []string
	Output  string // "text" (default)
}

// ParseFlags parses command-line arguments for the stats subcommand.
// It expects args to be os.Args[2:] (after the subcommand name).
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)

	output := fs.String("output", "text", "Output format: text")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	sources := fs.Args()
	if len(sources) == 0 {
		return nil, fmt.Errorf("at least one source .env file is required")
	}

	return &Flags{
		Sources: sources,
		Output:  *output,
	}, nil
}

// Run executes the stats command using the provided flags, writing results to w.
func Run(f *Flags, w io.Writer) error {
	for _, src := range f.Sources {
		result, err := Analyze(src)
		if err != nil {
			return fmt.Errorf("analyze %s: %w", src, err)
		}

		fmt.Fprintf(w, "=== %s ===\n", src)
		WriteText(w, result)
		fmt.Fprintln(w)
	}
	return nil
}

// RunCLI is the top-level entry point for the stats subcommand.
// It parses os.Args, runs the command, and exits on error.
func RunCLI(args []string) {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff stats: %v\n", err)
		os.Exit(1)
	}

	if err := Run(f, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "envdiff stats: %v\n", err)
		os.Exit(1)
	}
}
