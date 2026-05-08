package copy

import (
	"flag"
	"fmt"
	"io"
)

// Flags holds parsed CLI arguments for the copy command.
type Flags struct {
	Source    string
	Dest      string
	Key       string
	NewKey    string
	Overwrite bool
}

// ParseFlags parses copy command flags from args.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("copy", flag.ContinueOnError)
	source := fs.String("source", "", "source .env file")
	dest := fs.String("dest", "", "destination .env file")
	key := fs.String("key", "", "key to copy from source")
	newKey := fs.String("new-key", "", "renamed key in destination (optional)")
	overwrite := fs.Bool("overwrite", false, "overwrite existing key in destination")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *source == "" {
		return Flags{}, fmt.Errorf("--source is required")
	}
	if *dest == "" {
		return Flags{}, fmt.Errorf("--dest is required")
	}
	if *key == "" {
		return Flags{}, fmt.Errorf("--key is required")
	}
	return Flags{
		Source:    *source,
		Dest:      *dest,
		Key:       *key,
		NewKey:    *newKey,
		Overwrite: *overwrite,
	}, nil
}

// Run executes the copy command, writing output to w.
func Run(args []string, w io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	res, err := CopyKey(Options{
		SourceFile: flags.Source,
		DestFile:   flags.Dest,
		Key:        flags.Key,
		NewKey:     flags.NewKey,
		Overwrite:  flags.Overwrite,
	})
	if err != nil {
		return err
	}

	destKey := res.NewKey
	action := "added"
	if res.Updated {
		action = "updated"
	}
	fmt.Fprintf(w, "key %q %s in %s\n", destKey, action, flags.Dest)
	return nil
}
