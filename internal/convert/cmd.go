package convert

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// Flags holds parsed CLI flags for the convert command.
type Flags struct {
	Source string
	Format string
	Output string
}

// ParseFlags parses convert subcommand arguments.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	source := fs.String("source", "", "source .env file (required)")
	format := fs.String("format", "dotenv", "output format: dotenv, export, json, yaml")
	output := fs.String("output", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *source == "" {
		return Flags{}, fmt.Errorf("--source is required")
	}
	return Flags{Source: *source, Format: *format, Output: *output}, nil
}

// Run executes the convert command with the given flags, writing to w.
func Run(f Flags, w io.Writer) error {
	env, err := parser.ParseFile(f.Source)
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}

	res, err := Convert(env, Format(f.Format))
	if err != nil {
		return err
	}

	if f.Output != "" {
		if err := os.WriteFile(f.Output, []byte(res.Output), 0644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		fmt.Fprintf(w, "Converted %s → %s (%s format)\n", f.Source, f.Output, f.Format)
		return nil
	}

	_, err = fmt.Fprint(w, res.Output)
	return err
}
