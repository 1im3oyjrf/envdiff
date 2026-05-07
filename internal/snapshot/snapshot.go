package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/parser"
)

// Snapshot represents a captured state of an .env file at a point in time.
type Snapshot struct {
	File      string            `json:"file"`
	CapturedAt time.Time        `json:"captured_at"`
	Entries   map[string]string `json:"entries"`
}

// Capture reads the given .env file and returns a Snapshot of its current state.
func Capture(filePath string) (*Snapshot, error) {
	entries, err := parser.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to parse %q: %w", filePath, err)
	}
	return &Snapshot{
		File:       filePath,
		CapturedAt: time.Now().UTC(),
		Entries:    entries,
	}, nil
}

// Save writes the snapshot to a JSON file at the given destination path.
func Save(s *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: failed to create file %q: %w", dest, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: failed to encode snapshot: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to open %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: failed to decode snapshot: %w", err)
	}
	return &s, nil
}

// Diff compares two snapshots and returns keys that were added, removed, or changed.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
}

// Compare returns the difference between an old and new snapshot.
func Compare(old, new *Snapshot) DiffResult {
	result := DiffResult{}

	for k, newVal := range new.Entries {
		oldVal, exists := old.Entries[k]
		if !exists {
			result.Added = append(result.Added, k)
		} else if oldVal != newVal {
			result.Changed = append(result.Changed, k)
		}
	}

	for k := range old.Entries {
		if _, exists := new.Entries[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	return result
}
