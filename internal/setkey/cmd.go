package setkey

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the setkey command.
type Flags struct {
	Source    string
	Key       string
	Value     string
	Overwrite bool
}

// ParseFlags parses command-line arguments for the setkey subcommand.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("setkey", flag.ContinueOnError)
	source := fs.String("source", "", "path to the .env file")
	key := fs.String("key", "", "key to set")
	value := fs.String("value", "", "value to assign")
	overwrite := fs.Bool("overwrite", false, "overwrite key if it already exists")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *source == "" {
		return Flags{}, fmt.Errorf("--source is required")
	}
	if *key == "" {
		return Flags{}, fmt.Errorf("--key is required")
	}
	return Flags{
		Source:    *source,
		Key:       *key,
		Value:     *value,
		Overwrite: *overwrite,
	}, nil
}

// Run executes the setkey command, writing output to w.
func Run(args []string, w io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	res, err := SetKey(flags.Source, flags.Key, flags.Value, flags.Overwrite)
	if err != nil {
		return err
	}

	action := "added"
	if res.Updated {
		action = "updated"
	}
	fmt.Fprintf(w, "%s %s=%s in %s\n", action, res.Key, res.Value, flags.Source)
	return nil
}

// RunCLI is the entry point wired into main.
func RunCLI(args []string) {
	if err := Run(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
