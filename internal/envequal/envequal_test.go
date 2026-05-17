package envequal_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envdiff/internal/envequal"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCheckEqual_Identical(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	r, err := envequal.CheckEqual(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Equal {
		t.Errorf("expected files to be equal")
	}
	if len(r.OnlyInA)+len(r.OnlyInB)+len(r.Mismatched) != 0 {
		t.Errorf("expected no differences")
	}
}

func TestCheckEqual_OnlyInA(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\nEXTRA=yes\n")
	b := writeTempEnv(t, "FOO=bar\n")
	r, err := envequal.CheckEqual(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Equal {
		t.Errorf("expected files to differ")
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA only in A, got %v", r.OnlyInA)
	}
}

func TestCheckEqual_OnlyInB(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=bar\nNEW=val\n")
	r, err := envequal.CheckEqual(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Equal {
		t.Errorf("expected files to differ")
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "NEW" {
		t.Errorf("expected NEW only in B, got %v", r.OnlyInB)
	}
}

func TestCheckEqual_Mismatched(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=baz\n")
	r, err := envequal.CheckEqual(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Equal {
		t.Errorf("expected files to differ")
	}
	if len(r.Mismatched) != 1 || r.Mismatched[0] != "FOO" {
		t.Errorf("expected FOO mismatch, got %v", r.Mismatched)
	}
}

func TestCheckEqual_MissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.env")
	_, err := envequal.CheckEqual(missing, missing)
	if err == nil {
		t.Errorf("expected error for missing file")
	}
}

func TestWriteText_Equal(t *testing.T) {
	r := envequal.Result{Equal: true, FileA: "a.env", FileB: "b.env"}
	var buf bytes.Buffer
	envequal.WriteText(&buf, r)
	if !strings.Contains(buf.String(), "identical") {
		t.Errorf("expected 'identical' in output, got: %s", buf.String())
	}
}

func TestWriteText_Differ(t *testing.T) {
	r := envequal.Result{
		Equal:      false,
		FileA:      "a.env",
		FileB:      "b.env",
		Mismatched: []string{"SECRET_KEY"},
	}
	var buf bytes.Buffer
	envequal.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "differ") {
		t.Errorf("expected 'differ' in output")
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output")
	}
}
