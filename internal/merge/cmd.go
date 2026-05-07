package merge

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI flags for the merge subcommand.
type Flags struct {
	Source   string
	Target   string
	Output   string
	Strategy Strategy
}

// ParseFlags parses merge subcommand arguments.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	source := fs.String("source", "", "source .env file (required)")
	target := fs.String("target", "", "target .env file (required)")
	output := fs.String("output", "merged.env", "output file path")
	strategy := fs.String("strategy", "source", "conflict resolution strategy: source|target|union")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *source == "" {
		return nil, fmt.Errorf("--source is required")
	}
	if *target == "" {
		return nil, fmt.Errorf("--target is required")
	}

	var s Strategy
	switch *strategy {
	case "source":
		s = StrategySource
	case "target":
		s = StrategyTarget
	case "union":
		s = StrategyUnion
	default:
		return nil, fmt.Errorf("unknown strategy %q: use source, target, or union", *strategy)
	}

	return &Flags{Source: *source, Target: *target, Output: *output, Strategy: s}, nil
}

// Run executes the merge command, writing output and printing a summary.
func Run(args []string, out io.Writer) error {
	flags, err := ParseFlags(args)
	if err != nil {
		return err
	}

	result, err := Merge(flags.Source, flags.Target, flags.Strategy)
	if err != nil {
		return err
	}

	if err := Write(result, flags.Output); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	fmt.Fprintf(out, "Merged %d keys into %s\n", len(result.Entries), flags.Output)
	if len(result.Conflicts) > 0 {
		fmt.Fprintf(out, "%d conflict(s) resolved using strategy=%s:\n", len(result.Conflicts), flags.Strategy)
		for _, c := range result.Conflicts {
			fmt.Fprintf(out, "  %s: %q vs %q -> %q\n", c.Key, c.SourceValue, c.TargetValue, c.Resolved)
		}
	}

	if _, err := os.Stat(flags.Output); err == nil {
		fmt.Fprintf(out, "Output written to %s\n", flags.Output)
	}
	return nil
}
