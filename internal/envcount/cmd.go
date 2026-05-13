package envcount

import (
	"flag"
	"fmt"
	"os"
)

// Flags holds parsed CLI options for the envcount command.
type Flags struct {
	Sources []string
}

// ParseFlags parses command-line arguments for the envcount subcommand.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("envcount", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: envdiff envcount <file1.env> [file2.env ...]")
		fmt.Fprintln(fs.Output(), "\nCounts keys in one or more .env files.")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	sources := fs.Args()
	if len(sources) == 0 {
		return nil, fmt.Errorf("envcount: at least one source file is required")
	}

	return &Flags{Sources: sources}, nil
}

// Run executes the envcount command with the given flags.
func Run(f *Flags) error {
	result, err := CountFiles(f.Sources)
	if err != nil {
		return err
	}
	WriteText(os.Stdout, result)
	return nil
}

// RunCLI is the entry point called from main with raw os.Args.
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
