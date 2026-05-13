package envcount

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestCountFiles_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	r, err := CountFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.TotalFiles != 1 {
		t.Errorf("expected 1 file, got %d", r.TotalFiles)
	}
	if r.TotalKeys != 2 {
		t.Errorf("expected 2 keys, got %d", r.TotalKeys)
	}
	if r.Files[0].Set != 2 {
		t.Errorf("expected 2 set keys, got %d", r.Files[0].Set)
	}
}

func TestCountFiles_EmptyValues(t *testing.T) {
	path := writeTempEnv(t, "FOO=\nBAR=hello\nBAZ=\n")
	r, err := CountFiles([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Files[0].Empty != 2 {
		t.Errorf("expected 2 empty keys, got %d", r.Files[0].Empty)
	}
	if r.Files[0].Set != 1 {
		t.Errorf("expected 1 set key, got %d", r.Files[0].Set)
	}
}

func TestCountFiles_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\nB=2\n")
	p2 := writeTempEnv(t, "C=3\n")
	r, err := CountFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.TotalFiles != 2 {
		t.Errorf("expected 2 files, got %d", r.TotalFiles)
	}
	if r.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", r.TotalKeys)
	}
}

func TestCountFiles_MissingFile(t *testing.T) {
	_, err := CountFiles([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWriteText_Output(t *testing.T) {
	r := &Result{
		TotalFiles: 1,
		TotalKeys:  3,
		Files: []FileCount{
			{File: "prod.env", Total: 3, Set: 2, Empty: 1},
		},
	}
	var buf bytes.Buffer
	WriteText(&buf, r)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains([]byte(out), []byte("prod.env")) {
		t.Error("expected output to contain filename")
	}
	if !bytes.Contains([]byte(out), []byte("Total:")) {
		t.Error("expected output to contain 'Total:'")
	}
}
