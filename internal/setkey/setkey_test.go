package setkey

import (
	"os"
	"path/filepath"
	"strings"
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

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	return string(b)
}

func TestSetKey_AddNewKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := SetKey(path, "NEW_KEY", "hello", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Updated {
		t.Error("expected Updated=false for new key")
	}
	content := readFile(t, path)
	if !strings.Contains(content, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY=hello in file, got:\n%s", content)
	}
}

func TestSetKey_UpdateExistingKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nSECRET=old\n")
	res, err := SetKey(path, "SECRET", "new", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Updated {
		t.Error("expected Updated=true for existing key")
	}
	content := readFile(t, path)
	if !strings.Contains(content, "SECRET=new") {
		t.Errorf("expected SECRET=new in file, got:\n%s", content)
	}
	if strings.Contains(content, "SECRET=old") {
		t.Error("old value should have been replaced")
	}
}

func TestSetKey_NoOverwriteReturnsError(t *testing.T) {
	path := writeTempEnv(t, "EXISTING=value\n")
	_, err := SetKey(path, "EXISTING", "new", false)
	if err == nil {
		t.Error("expected error when overwrite=false and key exists")
	}
}

func TestSetKey_PreservesCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\nFOO=bar\n\nBAZ=qux\n")
	_, err := SetKey(path, "FOO", "updated", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readFile(t, path)
	if !strings.Contains(content, "# comment") {
		t.Error("comment should be preserved")
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Error("unrelated key should be preserved")
	}
}

func TestSetKey_CreatesFileIfNotExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.env")
	res, err := SetKey(path, "BRAND_NEW", "yes", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Updated {
		t.Error("expected Updated=false for brand new file")
	}
	content := readFile(t, path)
	if !strings.Contains(content, "BRAND_NEW=yes") {
		t.Errorf("expected BRAND_NEW=yes in new file, got:\n%s", content)
	}
}
