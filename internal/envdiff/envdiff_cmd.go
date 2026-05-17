package envdiff

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// CmdFlags holds parsed CLI flags for the envdiff subcommand.
type CmdFlags struct {
	FileA      string
	FileB      string
	MaskSecret bool
	JSON       bool
	Output     string
}

// ParseFlags parses CLI arguments for the envdiff subcommand.
func ParseFlags(args []string) (*CmdFlags, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	mask := fs.Bool("mask", false, "mask secret values in output")
	jsonOut := fs.Bool("json", false, "output results as JSON")
	output := fs.String("output", "", "write output to file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	positional := fs.Args()
	if len(positional) < 2 {
		return nil, fmt.Errorf("usage: envdiff [flags] <file-a> <file-b>")
	}

	return &CmdFlags{
		FileA:      positional[0],
		FileB:      positional[1],
		MaskSecret: *mask,
		JSON:       *jsonOut,
		Output:     *output,
	}, nil
}

// Run executes the envdiff command with the given flags.
func Run(flags *CmdFlags, stdout io.Writer) error {
	result, err := Diff(flags.FileA, flags.FileB)
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}

	var w io.Writer = stdout
	if flags.Output != "" {
		f, err := os.Create(flags.Output)
		if err != nil {
			return fmt.Errorf("cannot open output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	WriteText(w, result, flags.MaskSecret)
	return nil
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	flags, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := Run(flags, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
