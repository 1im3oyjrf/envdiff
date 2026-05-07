package promote

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the outcome of a promotion operation.
type Result struct {
	Promoted []string
	Skipped  []string
	Overwritten []string
}

// Options controls promote behaviour.
type Options struct {
	// Overwrite existing keys in target when true.
	Overwrite bool
	// Keys to explicitly include; if empty all missing keys are promoted.
	Keys []string
}

// Promote copies keys present in source but missing in target into the target
// map. It returns a Result describing what happened.
func Promote(sourceFile, targetFile string, opts Options) (map[string]string, Result, error) {
	src, err := parser.ParseFile(sourceFile)
	if err != nil {
		return nil, Result{}, fmt.Errorf("reading source: %w", err)
	}

	dst, err := parser.ParseFile(targetFile)
	if err != nil {
		return nil, Result{}, fmt.Errorf("reading target: %w", err)
	}

	filter := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		filter[k] = true
	}

	result := Result{}
	for k, v := range src {
		if len(filter) > 0 && !filter[k] {
			continue
		}
		if _, exists := dst[k]; exists {
			if !opts.Overwrite {
				result.Skipped = append(result.Skipped, k)
				continue
			}
			result.Overwritten = append(result.Overwritten, k)
		} else {
			result.Promoted = append(result.Promoted, k)
		}
		dst[k] = v
	}

	sort.Strings(result.Promoted)
	sort.Strings(result.Skipped)
	sort.Strings(result.Overwritten)

	return dst, result, nil
}
