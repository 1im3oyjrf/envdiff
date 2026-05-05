package diff

import "github.com/yourusername/envdiff/internal/parser"

// Result holds the comparison outcome between two env files.
type Result struct {
	// MissingInTarget contains keys present in source but absent in target.
	MissingInTarget []string
	// MissingInSource contains keys present in target but absent in source.
	MissingInSource []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchEntry
}

// MismatchEntry describes a key whose value differs between source and target.
type MismatchEntry struct {
	Key         string
	SourceValue string
	TargetValue string
}

// IsClean returns true when there are no differences between the two env files.
func (r Result) IsClean() bool {
	return len(r.MissingInTarget) == 0 &&
		len(r.MissingInSource) == 0 &&
		len(r.Mismatched) == 0
}

// Compare computes the diff between two EnvMaps.
// source is treated as the reference (e.g. .env.example).
// target is the environment being validated (e.g. .env.production).
func Compare(source, target parser.EnvMap) Result {
	var result Result

	for key, srcVal := range source {
		tgtVal, exists := target[key]
		if !exists {
			result.MissingInTarget = append(result.MissingInTarget, key)
			continue
		}
		if srcVal != tgtVal {
			result.Mismatched = append(result.Mismatched, MismatchEntry{
				Key:         key,
				SourceValue: srcVal,
				TargetValue: tgtVal,
			})
		}
	}

	for key := range target {
		if _, exists := source[key]; !exists {
			result.MissingInSource = append(result.MissingInSource, key)
		}
	}

	sortStrings(result.MissingInTarget)
	sortStrings(result.MissingInSource)
	sortMismatched(result.Mismatched)

	return result
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func sortMismatched(m []MismatchEntry) {
	for i := 1; i < len(m); i++ {
		for j := i; j > 0 && m[j].Key < m[j-1].Key; j-- {
			m[j], m[j-1] = m[j-1], m[j]
		}
	}
}
