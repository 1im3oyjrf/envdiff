package envkeys_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envkeys"
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

func TestListKeys_Basic(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	r, err := envkeys.ListKeys(path, envkeys.Options{Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	if r.Keys[0] != "BAZ" || r.Keys[1] != "FOO" {
		t.Errorf("unexpected order: %v", r.Keys)
	}
}

func TestListKeys_PrefixFilter(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	r, err := envkeys.ListKeys(path, envkeys.Options{Sorted: true, PrefixFilter: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	for _, k := range r.Keys {
		if !strings.HasPrefix(k, "DB_") {
			t.Errorf("key %q does not have expected prefix", k)
		}
	}
}

func TestListKeys_MissingFile(t *testing.T) {
	_, err := envkeys.ListKeys(filepath.Join(t.TempDir(), "nonexistent.env"), envkeys.Options{})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWriteText_KeysOnly(t *testing.T) {
	path := writeTempEnv(t, "ALPHA=1\nBETA=2\n")
	r, _ := envkeys.ListKeys(path, envkeys.Options{Sorted: true})
	var buf bytes.Buffer
	envkeys.WriteText(&buf, r, envkeys.Options{})
	out := buf.String()
	if !strings.Contains(out, "ALPHA") || !strings.Contains(out, "BETA") {
		t.Errorf("expected keys in output, got: %s", out)
	}
	if strings.Contains(out, "= 1") {
		t.Errorf("values should not appear in keys-only mode")
	}
}

func TestWriteText_WithValues(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080\n")
	r, _ := envkeys.ListKeys(path, envkeys.Options{Sorted: true})
	var buf bytes.Buffer
	envkeys.WriteText(&buf, r, envkeys.Options{ValuesOnly: true})
	out := buf.String()
	if !strings.Contains(out, "PORT") || !strings.Contains(out, "8080") {
		t.Errorf("expected key and value in output, got: %s", out)
	}
}
