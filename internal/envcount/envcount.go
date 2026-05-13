package envcount

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// FileCount holds the key count result for a single file.
type FileCount struct {
	File  string
	Total int
	Empty int
	Set   int
}

// Result holds the aggregated counts across all files.
type Result struct {
	Files      []FileCount
	TotalFiles int
	TotalKeys  int
}

// CountFiles parses each file and returns key counts per file.
func CountFiles(paths []string) (*Result, error) {
	result := &Result{
		TotalFiles: len(paths),
	}

	for _, path := range paths {
		env, err := parser.ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("envcount: failed to parse %s: %w", path, err)
		}

		fc := FileCount{
			File:  path,
			Total: len(env),
		}

		for _, v := range env {
			if v == "" {
				fc.Empty++
			} else {
				fc.Set++
			}
		}

		result.Files = append(result.Files, fc)
		result.TotalKeys += fc.Total
	}

	sort.Slice(result.Files, func(i, j int) bool {
		return result.Files[i].File < result.Files[j].File
	})

	return result, nil
}

// WriteText writes a human-readable summary of key counts to w.
func WriteText(w io.Writer, r *Result) {
	for _, fc := range r.Files {
		fmt.Fprintf(w, "%s: %d keys (%d set, %d empty)\n", fc.File, fc.Total, fc.Set, fc.Empty)
	}
	fmt.Fprintf(w, "\nTotal: %d files, %d keys\n", r.TotalFiles, r.TotalKeys)
}
