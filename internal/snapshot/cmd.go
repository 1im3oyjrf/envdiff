package snapshot

import (
	"flag"
	"fmt"
	"os"
)

// Flags holds parsed CLI options for the snapshot subcommand.
type Flags struct {
	Action     string // "capture" or "compare"
	EnvFile    string
	SnapshotFile string
	MaskSecrets bool
}

// ParseFlags parses snapshot subcommand arguments.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)

	envFile := fs.String("env", "", "Path to the .env file")
	snapshotFile := fs.String("snapshot", ".env.snapshot", "Path to the snapshot file")
	maskSecrets := fs.Bool("mask", false, "Mask secret values in output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() < 1 {
		return nil, fmt.Errorf("action required: capture or compare")
	}

	action := fs.Arg(0)
	if action != "capture" && action != "compare" {
		return nil, fmt.Errorf("unknown action %q: must be capture or compare", action)
	}

	if *envFile == "" {
		return nil, fmt.Errorf("--env flag is required")
	}

	return &Flags{
		Action:       action,
		EnvFile:      *envFile,
		SnapshotFile: *snapshotFile,
		MaskSecrets:  *maskSecrets,
	}, nil
}

// Run executes the snapshot action described by f, writing output to stdout.
func Run(f *Flags, out *os.File) error {
	switch f.Action {
	case "capture":
		snap, err := Capture(f.EnvFile)
		if err != nil {
			return fmt.Errorf("capture failed: %w", err)
		}
		if err := Save(snap, f.SnapshotFile); err != nil {
			return fmt.Errorf("save failed: %w", err)
		}
		fmt.Fprintf(out, "Snapshot saved to %s (%d keys)\n", f.SnapshotFile, len(snap.Entries))

	case "compare":
		current, err := Capture(f.EnvFile)
		if err != nil {
			return fmt.Errorf("capture failed: %w", err)
		}
		previous, err := Load(f.SnapshotFile)
		if err != nil {
			return fmt.Errorf("load snapshot failed: %w", err)
		}
		result := Compare(previous, current)
		WriteSummary(out, result, f.MaskSecrets)
	}
	return nil
}
