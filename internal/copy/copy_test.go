package copy

import (
	"os"
	"path/filepath"
	"strings"
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

func TestCopyKey_NewKey(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	dest := filepath.Join(t.TempDir(), "dest.env")

	res, err := CopyKey(Options{SourceFile: src, DestFile: dest, Key: "DB_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Created != true || res.Updated != false {
		t.Errorf("expected Created=true, got %+v", res)
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "DB_HOST=localhost") {
		t.Errorf("dest missing expected key, got: %s", data)
	}
}

func TestCopyKey_Rename(t *testing.T) {
	src := writeTempEnv(t, "API_SECRET=abc123\n")
	dest := writeTempEnv(t, "OTHER=val\n")

	res, err := CopyKey(Options{SourceFile: src, DestFile: dest, Key: "API_SECRET", NewKey: "SERVICE_SECRET"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.NewKey != "SERVICE_SECRET" {
		t.Errorf("expected NewKey=SERVICE_SECRET, got %s", res.NewKey)
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "SERVICE_SECRET=abc123") {
		t.Errorf("renamed key not found in dest: %s", data)
	}
}

func TestCopyKey_NoOverwriteReturnsError(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=prod-host\n")
	dest := writeTempEnv(t, "DB_HOST=localhost\n")

	_, err := CopyKey(Options{SourceFile: src, DestFile: dest, Key: "DB_HOST", Overwrite: false})
	if err == nil {
		t.Fatal("expected error when key exists and overwrite=false")
	}
}

func TestCopyKey_Overwrite(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=prod-host\n")
	dest := writeTempEnv(t, "DB_HOST=localhost\n")

	res, err := CopyKey(Options{SourceFile: src, DestFile: dest, Key: "DB_HOST", Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Updated {
		t.Error("expected Updated=true")
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "DB_HOST=prod-host") {
		t.Errorf("expected overwritten value, got: %s", data)
	}
}

func TestCopyKey_MissingSourceKey(t *testing.T) {
	src := writeTempEnv(t, "OTHER=val\n")
	dest := writeTempEnv(t, "")

	_, err := CopyKey(Options{SourceFile: src, DestFile: dest, Key: "MISSING_KEY"})
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestCopyKey_MissingSourceFile(t *testing.T) {
	_, err := CopyKey(Options{SourceFile: "/nonexistent.env", DestFile: "/tmp/x.env", Key: "K"})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}
