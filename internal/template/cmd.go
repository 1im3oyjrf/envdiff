package template

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// Options holds configuration for the template generation command.
type Options struct {
	Sources []string
	Output  string
}

// ParseFlags parses template subcommand flags from args.
func ParseFlags(args []string, stderr io.Writer) (*Options, error) {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var output string
	fs.StringVar(&output, "output", ".env.template", "path to write the generated template")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	sources := fs.Args()
	if len(sources) == 0 {
		return nil, fmt.Errorf("template: at least one source file is required")
	}

	return &Options{
		Sources: sources,
		Output:  output,
	}, nil
}

// Run executes the template generation command and writes output.
func Run(opts *Options, stdout io.Writer) error {
	entries, err := GenerateFromFiles(opts.Sources)
	if err != nil {
		return err
	}

	if err := Write(opts.Output, entries); err != nil {
		return err
	}

	keys := make([]string, len(entries))
	for i, e := range entries {
		keys[i] = e.Key
	}
	fmt.Fprintf(stdout, "Generated template with %d keys: %s\n",
		len(entries), strings.Join(keys, ", "))
	return nil
}
