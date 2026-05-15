package envclone

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI flags for the clone command.
type Flags struct {
	Source      string
	Dest        string
	Prefix      string
	StripPrefix bool
	Overwrite   bool
}

// ParseFlags parses command-line arguments for the clone sub-command.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	var f Flags
	fs.StringVar(&f.Source, "source", "", "source .env file (required)")
	fs.StringVar(&f.Dest, "dest", "", "destination .env file (required)")
	fs.StringVar(&f.Prefix, "prefix", "", "only clone keys with this prefix")
	fs.BoolVar(&f.StripPrefix, "strip-prefix", false, "strip prefix from cloned keys")
	fs.BoolVar(&f.Overwrite, "overwrite", false, "overwrite existing keys in destination")
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if f.Source == "" {
		return Flags{}, fmt.Errorf("--source is required")
	}
	if f.Dest == "" {
		return Flags{}, fmt.Errorf("--dest is required")
	}
	return f, nil
}

// Run executes the clone command, writing output to w.
func Run(args []string, w io.Writer) error {
	f, err := ParseFlags(args)
	if err != nil {
		return err
	}
	res, err := CloneFile(Options{
		Source:      f.Source,
		Dest:        f.Dest,
		Prefix:      f.Prefix,
		StripPrefix: f.StripPrefix,
		Overwrite:   f.Overwrite,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Cloned: %d key(s) copied, %d overwritten, %d skipped\n",
		len(res.Copied), len(res.Overwrite), len(res.Skipped))
	for _, k := range res.Copied {
		fmt.Fprintf(w, "  + %s\n", k)
	}
	for _, k := range res.Overwrite {
		fmt.Fprintf(w, "  ~ %s\n", k)
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(w, "  - %s (skipped)\n", k)
	}
	return nil
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) {
	if err := Run(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
