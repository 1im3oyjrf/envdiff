package template

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestGenerateFromFiles_Single(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	entries, err := GenerateFromFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "APP_NAME" || entries[1].Key != "DB_HOST" {
		t.Errorf("unexpected keys: %v", entries)
	}
}

func TestGenerateFromFiles_Deduplication(t *testing.T) {
	p1 := writeTempEnv(t, "APP_NAME=foo\nSECRET_KEY=abc\n")
	p2 := writeTempEnv(t, "APP_NAME=bar\nDB_URL=postgres\n")
	entries, err := GenerateFromFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestGenerateFromFiles_PreservesComments(t *testing.T) {
	path := writeTempEnv(t, "# database host\nDB_HOST=localhost\n")
	entries, err := GenerateFromFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Comment != "database host" {
		t.Errorf("expected comment 'database host', got %q", entries[0].Comment)
	}
}

func TestGenerateFromFiles_Sorted(t *testing.T) {
	path := writeTempEnv(t, "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")
	entries, err := GenerateFromFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "APPLE" || entries[1].Key != "MIDDLE" || entries[2].Key != "ZEBRA" {
		t.Errorf("entries not sorted: %v", entries)
	}
}

func TestWrite_CreatesFile(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Comment: "application name"},
		{Key: "DB_HOST"},
	}
	out := filepath.Join(t.TempDir(), "template.env")
	if err := Write(out, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output: %v", err)
	}
	content := string(data)
	if content != "# application name\nAPP_NAME=\nDB_HOST=\n" {
		t.Errorf("unexpected content:\n%s", content)
	}
}

func TestGenerateFromFiles_MissingFile(t *testing.T) {
	_, err := GenerateFromFiles([]string{"/nonexistent/path.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
