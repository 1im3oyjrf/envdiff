package envequal

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds parsed CLI options for the envequal command.
type Flags struct {
	FileA  string
	FileB  string
	Quiet  bool
}

// ParseFlags parses command-line arguments for the envequal subcommand.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("envequal", flag.ContinueOnError)
	quiet := fs.Bool("quiet", false, "suppress output; use exit code only")
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	if fs.NArg() < 2 {
		return Flags{}, fmt.Errorf("usage: envequal [--quiet] <file-a> <file-b>")
	}
	return Flags{
		FileA: fs.Arg(0),
		FileB: fs.Arg(1),
		Quiet: *quiet,
	}, nil
}

// Run executes the envequal command, writing results to w.
// Returns true if the files are equal.
func Run(f Flags, w io.Writer) (bool, error) {
	r, err := CheckEqual(f.FileA, f.FileB)
	if err != nil {
		return false, err
	}
	if !f.Quiet {
		WriteText(w, r)
	}
	return r.Equal, nil
}

// RunCLI is the entry point called from main.
func RunCLI(args []string) int {
	f, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	equal, err := Run(f, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	if !equal {
		return 1
	}
	return 0
}
