package promote

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/merge"
)

// Flags holds parsed CLI flags for the promote command.
type Flags struct {
	Source    string
	Target    string
	Output    string
	Overwrite bool
	Keys      []string
}

// ParseFlags parses promote sub-command flags from args.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	var (
		source    = fs.String("source", "", "source .env file")
		target    = fs.String("target", "", "target .env file")
		output    = fs.String("output", "", "output file (default: stdout)")
		overwrite = fs.Bool("overwrite", false, "overwrite existing keys in target")
		keys      = fs.String("keys", "", "comma-separated list of keys to promote")
	)
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if *source == "" {
		return Flags{}, fmt.Errorf("--source is required")
	}
	if *target == "" {
		return Flags{}, fmt.Errorf("--target is required")
	}
	var keyList []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keyList = append(keyList, k)
			}
		}
	}
	return Flags{
		Source:    *source,
		Target:    *target,
		Output:    *output,
		Overwrite: *overwrite,
		Keys:      keyList,
	}, nil
}

// Run executes the promote command writing results to w.
func Run(args []string, w io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	merged, result, err := Promote(flags.Source, flags.Target, Options{
		Overwrite: flags.Overwrite,
		Keys:      flags.Keys,
	})
	if err != nil {
		return err
	}

	var out io.Writer = w
	if flags.Output != "" {
		f, err := os.Create(flags.Output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	if err := merge.Write(merged, out); err != nil {
		return err
	}

	fmt.Fprintf(w, "promoted: %d  overwritten: %d  skipped: %d\n",
		len(result.Promoted), len(result.Overwritten), len(result.Skipped))
	return nil
}
