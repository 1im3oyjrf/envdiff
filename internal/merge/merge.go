package merge

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Strategy defines how conflicts are resolved during merge.
type Strategy string

const (
	StrategySource Strategy = "source" // prefer source values on conflict
	StrategyTarget Strategy = "target" // prefer target values on conflict
	StrategyUnion  Strategy = "union"  // include all keys from both files
)

// Result holds the merged key-value pairs and metadata.
type Result struct {
	Entries    map[string]string
	Conflicts  []Conflict
	Strategy   Strategy
}

// Conflict represents a key that existed in both files with different values.
type Conflict struct {
	Key         string
	SourceValue string
	TargetValue string
	Resolved    string
}

// Merge combines two .env files according to the given strategy.
func Merge(sourceFile, targetFile string, strategy Strategy) (*Result, error) {
	source, err := parser.ParseFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("reading source: %w", err)
	}
	target, err := parser.ParseFile(targetFile)
	if err != nil {
		return nil, fmt.Errorf("reading target: %w", err)
	}

	result := &Result{
		Entries:  make(map[string]string),
		Strategy: strategy,
	}

	// Add all source entries.
	for k, v := range source {
		result.Entries[k] = v
	}

	// Merge target entries, detecting conflicts.
	for k, tv := range target {
		if sv, exists := source[k]; exists && sv != tv {
			c := Conflict{Key: k, SourceValue: sv, TargetValue: tv}
			switch strategy {
			case StrategyTarget:
				c.Resolved = tv
				result.Entries[k] = tv
			default: // StrategySource or StrategyUnion
				c.Resolved = sv
			}
			result.Conflicts = append(result.Conflicts, c)
		} else if !exists {
			result.Entries[k] = tv
		}
	}

	return result, nil
}

// Write writes the merged result to the given path.
func Write(result *Result, outputPath string) error {
	keys := make([]string, 0, len(result.Entries))
	for k := range result.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, result.Entries[k])
	}

	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
