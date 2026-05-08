package redact

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the redact command.
type Flags struct {
	Source string
	Output string
}

// ParseFlags parses redact subcommand arguments.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("redact", flag.ContinueOnError)
	out := fs.String("output", "", "output file path (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if fs.NArg() < 1 {
		return Flags{}, fmt.Errorf("usage: envdiff redact <source> [--output <file>]")
	}
	return Flags{Source: fs.Arg(0), Output: *out}, nil
}

// Run executes the redact command using the provided arguments.
func Run(args []string, stdout io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	var dst io.Writer = stdout
	if flags.Output != "" {
		f, err := os.Create(flags.Output)
		if err != nil {
			return fmt.Errorf("redact: create output file: %w", err)
		}
		defer f.Close()
		dst = f
	}

	res, err := RedactFile(flags.Source, dst)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Redacted %d/%d lines in %s\n", res.LinesRedacted, res.LinesTotal, res.File)
	return nil
}
