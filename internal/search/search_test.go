package search_test

import (
	"os"
	"testing"

	"github.com/envdiff/internal/search"
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

func TestSearchFiles_KeyPattern(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")

	results, err := search.SearchFiles([]string{path}, search.Options{KeyPattern: "db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchFiles_ExactKey(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	results, err := search.SearchFiles([]string{path}, search.Options{KeyPattern: "DB_HOST", ExactKey: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected exactly DB_HOST, got %+v", results)
	}
}

func TestSearchFiles_ValuePattern(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nREDIS_HOST=localhost\nAPP_ENV=production\n")

	results, err := search.SearchFiles([]string{path}, search.Options{ValuePattern: "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchFiles_KeyAndValue(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nREDIS_HOST=localhost\n")

	results, err := search.SearchFiles([]string{path}, search.Options{KeyPattern: "db", ValuePattern: "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected only DB_HOST, got %+v", results)
	}
}

func TestSearchFiles_NoPattern(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")

	_, err := search.SearchFiles([]string{path}, search.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestSearchFiles_MissingFile(t *testing.T) {
	_, err := search.SearchFiles([]string{"/nonexistent/.env"}, search.Options{KeyPattern: "KEY"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSearchFiles_LineNumbers(t *testing.T) {
	path := writeTempEnv(t, "# comment\nFIRST=one\nSECOND=two\n")

	results, err := search.SearchFiles([]string{path}, search.Options{KeyPattern: "SECOND", ExactKey: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Line != 3 {
		t.Errorf("expected line 3, got %d", results[0].Line)
	}
}
