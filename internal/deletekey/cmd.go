package deletekey

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the deletekey command.
type Flags struct {
	File   string
	Key    string
	DryRun bool
}

// ParseFlags parses CLI arguments for the deletekey subcommand.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("deletekey", flag.ContinueOnError)
	file := fs.String("file", "", "path to the .env file (required)")
	key := fs.String("key", "", "key to delete (required)")
	dryRun := fs.Bool("dry-run", false, "preview deletion without modifying the file")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *file == "" {
		return Flags{}, fmt.Errorf("deletekey: --file is required")
	}
	if *key == "" {
		return Flags{}, fmt.Errorf("deletekey: --key is required")
	}
	return Flags{File: *file, Key: *key, DryRun: *dryRun}, nil
}

// Run executes the deletekey command, writing output to w.
func Run(f Flags, w io.Writer) error {
	res, err := DeleteKey(Options{
		File:   f.File,
		Key:    f.Key,
		DryRun: f.DryRun,
	})
	if err != nil {
		return err
	}

	if !res.Deleted {
		fmt.Fprintf(w, "key %q not found in %s\n", res.Key, res.File)
		return nil
	}

	if res.DryRun {
		fmt.Fprintf(w, "[dry-run] would delete key %q from %s\n", res.Key, res.File)
	} else {
		fmt.Fprintf(w, "deleted key %q from %s\n", res.Key, res.File)
	}
	return nil
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := Run(f, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
