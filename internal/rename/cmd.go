package rename

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI arguments for the rename command.
type Flags struct {
	File      string
	OldKey    string
	NewKey    string
	DryRun    bool
	Overwrite bool
}

// ParseFlags parses rename subcommand flags from args.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	file := fs.String("file", "", "path to .env file (required)")
	old := fs.String("old", "", "key to rename (required)")
	new_ := fs.String("new", "", "replacement key name (required)")
	dry := fs.Bool("dry-run", false, "preview changes without writing")
	over := fs.Bool("overwrite", false, "overwrite new key if it already exists")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *file == "" {
		return Flags{}, fmt.Errorf("--file is required")
	}
	if *old == "" {
		return Flags{}, fmt.Errorf("--old is required")
	}
	if *new_ == "" {
		return Flags{}, fmt.Errorf("--new is required")
	}
	return Flags{File: *file, OldKey: *old, NewKey: *new_, DryRun: *dry, Overwrite: *over}, nil
}

// Run executes the rename command, writing output to w.
func Run(args []string, w io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	opts := Options{DryRun: flags.DryRun, Overwrite: flags.Overwrite}
	res, err := RenameKey(flags.File, flags.OldKey, flags.NewKey, opts)
	if err != nil {
		return err
	}

	switch {
	case res.Renamed && flags.DryRun:
		fmt.Fprintf(w, "[dry-run] would rename %s → %s in %s\n", res.OldKey, res.NewKey, flags.File)
	case res.Renamed:
		fmt.Fprintf(w, "renamed %s → %s in %s\n", res.OldKey, res.NewKey, flags.File)
	case res.Skipped:
		fmt.Fprintf(w, "skipped: %s (%s)\n", res.OldKey, res.Reason)
	}

	if !res.Renamed && !res.Skipped {
		os.Exit(1)
	}
	return nil
}
