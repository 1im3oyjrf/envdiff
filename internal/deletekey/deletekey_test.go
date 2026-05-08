package deletekey

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

func TestDeleteKey_Success(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nSECRET=xyz\nBAZ=qux\n")
	res, err := DeleteKey(Options{File: p, Key: "SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Deleted {
		t.Error("expected Deleted=true")
	}
	contents := readFile(t, p)
	if strings.Contains(contents, "SECRET") {
		t.Errorf("key SECRET still present in file: %q", contents)
	}
	if !strings.Contains(contents, "FOO=bar") {
		t.Errorf("FOO=bar should remain: %q", contents)
	}
}

func TestDeleteKey_NotFound(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := DeleteKey(Options{File: p, Key: "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Deleted {
		t.Error("expected Deleted=false for missing key")
	}
}

func TestDeleteKey_DryRun(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nSECRET=xyz\n")
	res, err := DeleteKey(Options{File: p, Key: "SECRET", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Deleted {
		t.Error("expected Deleted=true even in dry-run")
	}
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
	contents := readFile(t, p)
	if !strings.Contains(contents, "SECRET=xyz") {
		t.Errorf("dry-run should not modify file; got: %q", contents)
	}
}

func TestDeleteKey_MissingFile(t *testing.T) {
	_, err := DeleteKey(Options{File: "/nonexistent/.env", Key: "FOO"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDeleteKey_PreservesCommentsAndBlanks(t *testing.T) {
	p := writeTempEnv(t, "# comment\nFOO=bar\n\nDEL=me\nBAZ=qux\n")
	_, err := DeleteKey(Options{File: p, Key: "DEL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	contents := readFile(t, p)
	if !strings.Contains(contents, "# comment") {
		t.Errorf("comment should be preserved: %q", contents)
	}
	if !strings.Contains(contents, "FOO=bar") {
		t.Errorf("FOO=bar should be preserved: %q", contents)
	}
	if strings.Contains(contents, "DEL=") {
		t.Errorf("DEL key should be removed: %q", contents)
	}
}
