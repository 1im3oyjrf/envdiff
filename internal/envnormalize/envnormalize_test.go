package envnormalize

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestNormalizeFile_TrimWhitespace(t *testing.T) {
	p := writeTempEnv(t, "  KEY = value  \nOTHER=ok\n")
	opts := DefaultOptions()
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Normalized != "KEY=value" {
		t.Errorf("expected 'KEY=value', got %q", results[0].Normalized)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true for trimmed line")
	}
	if results[1].Changed {
		t.Error("expected Changed=false for already clean line")
	}
}

func TestNormalizeFile_RemoveExport(t *testing.T) {
	p := writeTempEnv(t, "export SECRET=abc\nKEY=val\n")
	opts := DefaultOptions()
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Normalized != "SECRET=abc" {
		t.Errorf("expected 'SECRET=abc', got %q", results[0].Normalized)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true when export prefix removed")
	}
}

func TestNormalizeFile_LowercaseKeys(t *testing.T) {
	p := writeTempEnv(t, "DB_HOST=localhost\n")
	opts := DefaultOptions()
	opts.LowercaseKeys = true
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Normalized != "db_host=localhost" {
		t.Errorf("expected 'db_host=localhost', got %q", results[0].Normalized)
	}
}

func TestNormalizeFile_QuoteValues(t *testing.T) {
	p := writeTempEnv(t, "API_URL=https://example.com\n")
	opts := DefaultOptions()
	opts.QuoteValues = true
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `API_URL="https://example.com"`
	if results[0].Normalized != expected {
		t.Errorf("expected %q, got %q", expected, results[0].Normalized)
	}
}

func TestNormalizeFile_PreservesCommentsAndBlanks(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nKEY=val\n")
	opts := DefaultOptions()
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Normalized != "# comment" {
		t.Errorf("expected comment preserved, got %q", results[0].Normalized)
	}
	if results[1].Normalized != "" {
		t.Errorf("expected blank line preserved, got %q", results[1].Normalized)
	}
}

func TestNormalizeFile_MissingFile(t *testing.T) {
	_, err := NormalizeFile("/nonexistent/.env", DefaultOptions())
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWrite_RoundTrip(t *testing.T) {
	p := writeTempEnv(t, "export KEY=value\n")
	opts := DefaultOptions()
	results, err := NormalizeFile(p, opts)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	out := filepath.Join(t.TempDir(), "out.env")
	if err := Write(out, results); err != nil {
		t.Fatalf("write: %v", err)
	}
	data, _ := os.ReadFile(out)
	if string(data) != "KEY=value\n" {
		t.Errorf("unexpected output: %q", string(data))
	}
}
