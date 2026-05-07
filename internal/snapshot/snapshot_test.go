package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/snapshot"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCapture_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	s, err := snapshot.Capture(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Entries["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", s.Entries["FOO"])
	}
	if s.Entries["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", s.Entries["BAZ"])
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestCapture_MissingFile(t *testing.T) {
	_, err := snapshot.Capture("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nSECRET=abc123\n")
	s, err := snapshot.Capture(path)
	if err != nil {
		t.Fatalf("capture failed: %v", err)
	}
	s.CapturedAt = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, dest); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Entries["KEY"])
	}
	if !loaded.CapturedAt.Equal(s.CapturedAt) {
		t.Errorf("CapturedAt mismatch: got %v", loaded.CapturedAt)
	}
}

func TestCompare_Differences(t *testing.T) {
	old := &snapshot.Snapshot{
		Entries: map[string]string{"FOO": "1", "BAR": "2", "GONE": "old"},
	}
	new := &snapshot.Snapshot{
		Entries: map[string]string{"FOO": "1", "BAR": "changed", "NEW": "here"},
	}

	result := snapshot.Compare(old, new)

	if len(result.Added) != 1 || result.Added[0] != "NEW" {
		t.Errorf("expected Added=[NEW], got %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0] != "GONE" {
		t.Errorf("expected Removed=[GONE], got %v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed[0] != "BAR" {
		t.Errorf("expected Changed=[BAR], got %v", result.Changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	old := &snapshot.Snapshot{Entries: map[string]string{"X": "1"}}
	new := &snapshot.Snapshot{Entries: map[string]string{"X": "1"}}
	result := snapshot.Compare(old, new)
	if len(result.Added)+len(result.Removed)+len(result.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", result)
	}
}
