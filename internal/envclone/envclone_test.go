package envclone

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

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestCloneFile_BasicCopy(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dest := filepath.Join(t.TempDir(), "dest.env")

	res, err := CloneFile(Options{Source: src, Dest: dest})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	content := readFile(t, dest)
	if !strings.Contains(content, "FOO=bar") || !strings.Contains(content, "BAZ=qux") {
		t.Errorf("dest missing expected keys: %s", content)
	}
}

func TestCloneFile_SkipsExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=newval\nBAR=baz\n")
	dest := writeTempEnv(t, "FOO=oldval\n")

	res, err := CloneFile(Options{Source: src, Dest: dest, Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO skipped, got %v", res.Skipped)
	}
	if !strings.Contains(readFile(t, dest), "FOO=oldval") {
		t.Error("FOO should retain old value")
	}
}

func TestCloneFile_OverwriteExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=newval\n")
	dest := writeTempEnv(t, "FOO=oldval\n")

	res, err := CloneFile(Options{Source: src, Dest: dest, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwrite) != 1 {
		t.Errorf("expected 1 overwritten, got %v", res.Overwrite)
	}
	if !strings.Contains(readFile(t, dest), "FOO=newval") {
		t.Error("FOO should have new value")
	}
}

func TestCloneFile_PrefixFilter(t *testing.T) {
	src := writeTempEnv(t, "APP_FOO=1\nAPP_BAR=2\nDB_HOST=localhost\n")
	dest := filepath.Join(t.TempDir(), "dest.env")

	res, err := CloneFile(Options{Source: src, Dest: dest, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d: %v", len(res.Copied), res.Copied)
	}
	content := readFile(t, dest)
	if strings.Contains(content, "DB_HOST") {
		t.Error("DB_HOST should have been filtered out")
	}
}

func TestCloneFile_StripPrefix(t *testing.T) {
	src := writeTempEnv(t, "APP_FOO=hello\nAPP_BAR=world\n")
	dest := filepath.Join(t.TempDir(), "dest.env")

	_, err := CloneFile(Options{Source: src, Dest: dest, Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := readFile(t, dest)
	if !strings.Contains(content, "FOO=hello") {
		t.Errorf("expected FOO=hello after strip, got: %s", content)
	}
	if strings.Contains(content, "APP_FOO") {
		t.Error("prefix should have been stripped")
	}
}

func TestCloneFile_MissingSource(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "dest.env")
	_, err := CloneFile(Options{Source: "/nonexistent.env", Dest: dest})
	if err == nil {
		t.Error("expected error for missing source")
	}
}
