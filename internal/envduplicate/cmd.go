package envduplicate

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the envduplicate command.
type Flags struct {
	Files   []string
	Summary bool
	Output  string
}

// ParseFlags parses command-line arguments for the envduplicate subcommand.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("envduplicate", flag.ContinueOnError)
	summary := fs.Bool("summary", false, "print one-line summary per file")
	output := fs.String("output", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() == 0 {
		return nil, fmt.Errorf("at least one .env file is required")
	}

	return &Flags{
		Files:   fs.Args(),
		Summary: *summary,
		Output:  *output,
	}, nil
}

// Run executes the duplicate-value scan and writes results.
func Run(f *Flags, stdout io.Writer) error {
	results, err := FindDuplicateValues(f.Files)
	if err != nil {
		return err
	}

	w := stdout
	if f.Output != "" {
		file, err := os.Create(f.Output)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer file.Close()
		w = file
	}

	if f.Summary {
		WriteSummary(w, results)
	} else {
		WriteText(w, results)
	}

	return nil
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "envduplicate:", err)
		os.Exit(1)
	}
	if err := Run(f, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "envduplicate:", err)
		os.Exit(1)
	}
}
