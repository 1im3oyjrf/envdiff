package envhide

import (
	"flag"
	"fmt"
	"os"
)

// Flags holds parsed CLI options for the envhide command.
type Flags struct {
	Source      string
	Output      string
	Placeholder string
}

// ParseFlags parses command-line arguments for the envhide subcommand.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("envhide", flag.ContinueOnError)
	src := fs.String("source", "", "path to the .env file (required)")
	out := fs.String("output", "", "output file path (default: stdout)")
	ph := fs.String("placeholder", "***", "placeholder string for hidden values")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *src == "" {
		return Flags{}, fmt.Errorf("envhide: -source is required")
	}
	return Flags{Source: *src, Output: *out, Placeholder: *ph}, nil
}

// Run executes the envhide command using the provided flags.
func Run(f Flags) error {
	res, err := HideSecrets(f.Source, f.Placeholder)
	if err != nil {
		return err
	}
	if err := Write(res, f.Output); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "envhide: %d secret(s) hidden in %q\n", res.HiddenCount, res.File)
	return nil
}

// RunCLI is the entry point for the envhide subcommand.
func RunCLI(args []string) {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := Run(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
