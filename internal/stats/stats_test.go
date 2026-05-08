package stats

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-stats-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestAnalyze_BasicCounts(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\nDB_PASS=\n")
	fs, err := Analyze(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", fs.TotalKeys)
	}
	if fs.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", fs.EmptyValues)
	}
}

func TestAnalyze_SecretKeys(t *testing.T) {
	path := writeTempEnv(t, "API_KEY=abc123\nSECRET_TOKEN=xyz\nAPP_NAME=myapp\n")
	fs, err := Analyze(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.SecretKeys != 2 {
		t.Errorf("expected 2 secret keys, got %d", fs.SecretKeys)
	}
}

func TestAnalyze_UniqueValues(t *testing.T) {
	path := writeTempEnv(t, "A=foo\nB=foo\nC=bar\n")
	fs, err := Analyze(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", fs.UniqueValues)
	}
}

func TestAnalyze_MissingFile(t *testing.T) {
	_, err := Analyze("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWriteText_Output(t *testing.T) {
	fs := FileStats{
		FilePath:     ".env.production",
		TotalKeys:    5,
		EmptyValues:  1,
		SecretKeys:   2,
		UniqueValues: 4,
		DuplicateKeys: []string{},
	}
	var buf bytes.Buffer
	WriteText(&buf, fs)
	out := buf.String()

	for _, want := range []string{
		".env.production",
		"Total keys",
		"Secret keys",
		"Duplicate keys : none",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestWriteText_WithDuplicates(t *testing.T) {
	fs := FileStats{
		FilePath:      ".env",
		DuplicateKeys: []string{"DB_HOST", "PORT"},
	}
	var buf bytes.Buffer
	WriteText(&buf, fs)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected duplicate key DB_HOST in output, got:\n%s", out)
	}
}
