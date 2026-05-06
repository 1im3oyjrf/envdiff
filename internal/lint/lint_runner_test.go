package lint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func writeTempLintFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestRunFile_Clean(t *testing.T) {
	path := writeTempLintFile(t, "KEY=value\nANOTHER=123\n")
	result, err := lint.RunFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(result.Issues))
	}
	if result.File != path {
		t.Errorf("expected file %q, got %q", path, result.File)
	}
}

func TestRunFile_WithIssues(t *testing.T) {
	path := writeTempLintFile(t, "KEY=value\nDUPLICATE=1\nDUPLICATE=2\nEMPTY=\n")
	result, err := lint.RunFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) == 0 {
		t.Error("expected issues, got none")
	}
}

func TestRunFiles_MultipleFiles(t *testing.T) {
	path1 := writeTempLintFile(t, "KEY=value\n")
	path2 := writeTempLintFile(t, "ANOTHER=123\nDUP=1\nDUP=2\n")

	results, err := lint.RunFiles([]string{path1, path2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if len(results[0].Issues) != 0 {
		t.Errorf("expected no issues in first file, got %d", len(results[0].Issues))
	}
	if len(results[1].Issues) == 0 {
		t.Error("expected issues in second file, got none")
	}
}

func TestRunFiles_MissingFile(t *testing.T) {
	_, err := lint.RunFiles([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunFiles_Empty(t *testing.T) {
	results, err := lint.RunFiles([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
