package envimport

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Source represents a supported import source format.
type Source string

const (
	SourceDotenv Source = "dotenv"
	SourceExport Source = "export"
	SourceDocker  Source = "docker" // KEY=VALUE lines from docker inspect output
)

// Entry holds a single imported key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Options controls import behaviour.
type Options struct {
	Source     Source
	OutputFile string // if empty, write to stdout
	Overwrite  bool
}

// Import reads entries from srcPath according to opts.Source and writes
// them in standard dotenv format to opts.OutputFile (or stdout).
func Import(srcPath string, opts Options) ([]Entry, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("envimport: open %s: %w", srcPath, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		entry, ok := parseLine(line, opts.Source)
		if !ok {
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envimport: scan %s: %w", srcPath, err)
	}

	var out *os.File
	if opts.OutputFile == "" {
		out = os.Stdout
	} else {
		flags := os.O_CREATE | os.O_WRONLY
		if opts.Overwrite {
			flags |= os.O_TRUNC
		} else {
			flags |= os.O_APPEND
		}
		out, err = os.OpenFile(opts.OutputFile, flags, 0644)
		if err != nil {
			return nil, fmt.Errorf("envimport: open output %s: %w", opts.OutputFile, err)
		}
		defer out.Close()
	}

	w := bufio.NewWriter(out)
	for _, e := range entries {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return entries, w.Flush()
}

func parseLine(line string, src Source) (Entry, bool) {
	switch src {
	case SourceExport:
		line = strings.TrimPrefix(line, "export ")
		line = strings.TrimPrefix(line, "export\t")
	}
	idx := strings.IndexByte(line, '=')
	if idx <= 0 {
		return Entry{}, false
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	} else if len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'' {
		val = val[1 : len(val)-1]
	}
	if key == "" {
		return Entry{}, false
	}
	return Entry{Key: key, Value: val}, true
}
