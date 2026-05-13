package envimport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestImport_Dotenv(t *testing.T) {
	src := writeTempFile(t, "KEY1=value1\nKEY2=\"quoted\"\n# comment\n\nKEY3=simple\n")
	entries, err := Import(src, Options{Source: SourceDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Value != "quoted" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestImport_ExportFormat(t *testing.T) {
	src := writeTempFile(t, "export FOO=bar\nexport BAR='baz'\n")
	entries, err := Import(src, Options{Source: SourceExport})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[1].Value != "baz" {
		t.Errorf("expected unquoted single-quoted value, got %q", entries[1].Value)
	}
}

func TestImport_WritesToOutputFile(t *testing.T) {
	src := writeTempFile(t, "A=1\nB=2\n")
	outPath := filepath.Join(t.TempDir(), "out.env")
	_, err := Import(src, Options{Source: SourceDotenv, OutputFile: outPath, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "A=1") || !strings.Contains(string(data), "B=2") {
		t.Errorf("output missing expected keys: %s", data)
	}
}

func TestImport_MissingSourceFile(t *testing.T) {
	_, err := Import("/nonexistent/.env", Options{Source: SourceDotenv})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestImport_SkipsInvalidLines(t *testing.T) {
	src := writeTempFile(t, "VALID=yes\nNOEQUALS\n=nokey\nANOTHER=ok\n")
	entries, err := Import(src, Options{Source: SourceDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 valid entries, got %d", len(entries))
	}
}
