package promote

import (
	"os"
	"path/filepath"
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
	return filepath.Clean(f.Name())
}

func TestPromote_MissingKeys(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempEnv(t, "FOO=bar\n")

	out, res, err := Promote(src, dst, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "BAZ" {
		t.Errorf("expected BAZ promoted, got %v", res.Promoted)
	}
	if out["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux in output")
	}
}

func TestPromote_SkipsExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	_, res, err := Promote(src, dst, Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO skipped, got %v", res.Skipped)
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	out, res, err := Promote(src, dst, Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "FOO" {
		t.Errorf("expected FOO overwritten, got %v", res.Overwritten)
	}
	if out["FOO"] != "new" {
		t.Errorf("expected FOO=new after overwrite")
	}
}

func TestPromote_KeyFilter(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := writeTempEnv(t, "")

	_, res, err := Promote(src, dst, Options{Keys: []string{"BAR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "BAR" {
		t.Errorf("expected only BAR promoted, got %v", res.Promoted)
	}
}

func TestPromote_MissingSourceFile(t *testing.T) {
	dst := writeTempEnv(t, "FOO=bar\n")
	_, _, err := Promote("/nonexistent.env", dst, Options{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
