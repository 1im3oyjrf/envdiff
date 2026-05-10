package groupby

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the groupby command.
type Flags struct {
	Source string
}

// ParseFlags parses command-line arguments for the groupby subcommand.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("groupby", flag.ContinueOnError)
	source := fs.String("source", "", "path to .env file")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *source == "" {
		return nil, fmt.Errorf("--source is required")
	}
	return &Flags{Source: *source}, nil
}

// Run executes the groupby command, writing grouped output to w.
func Run(f *Flags, w io.Writer) error {
	groups, err := GroupByPrefix(f.Source)
	if err != nil {
		return err
	}
	WriteText(groups, w)
	return nil
}

// WriteText renders grouped env entries as human-readable text.
func WriteText(groups []Group, w io.Writer) {
	for _, g := range groups {
		fmt.Fprintf(w, "[%s] (%d keys)\n", g.Name, len(g.Entries))
		for _, e := range g.Entries {
			fmt.Fprintf(w, "  %s=%s\n", e.Key, e.Value)
		}
	}
	if len(groups) == 0 {
		fmt.Fprintln(w, "no entries found")
	}
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "groupby: %v\n", err)
		os.Exit(1)
	}
	if err := Run(f, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "groupby: %v\n", err)
		os.Exit(1)
	}
}
