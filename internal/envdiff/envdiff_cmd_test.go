package envdiff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvDiffFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseFlags_Valid(t *testing.T) {
	flags, err := ParseFlags([]string{"--mask", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.FileA != "a.env" || flags.FileB != "b.env" {
		t.Errorf("unexpected files: %v %v", flags.FileA, flags.FileB)
	}
	if !flags.MaskSecret {
		t.Error("expected MaskSecret to be true")
	}
}

func TestParseFlags_MissingArgs(t *testing.T) {
	_, err := ParseFlags([]string{"only-one.env"})
	if err == nil {
		t.Fatal("expected error for missing second file")
	}
}

func TestParseFlags_NoArgs(t *testing.T) {
	_, err := ParseFlags([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRun_WritesToStdout(t *testing.T) {
	a := writeTempEnvDiffFile(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnvDiffFile(t, "FOO=bar\n")

	flags := &CmdFlags{FileA: a, FileB: b}
	var buf bytes.Buffer
	if err := Run(flags, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "BAZ") {
		t.Errorf("expected BAZ in output, got: %s", out)
	}
}

func TestRun_WritesToFile(t *testing.T) {
	a := writeTempEnvDiffFile(t, "KEY=val\n")
	b := writeTempEnvDiffFile(t, "KEY=val\n")
	out := filepath.Join(t.TempDir(), "out.txt")

	flags := &CmdFlags{FileA: a, FileB: b, Output: out}
	var buf bytes.Buffer
	if err := Run(flags, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestRun_MissingFile(t *testing.T) {
	flags := &CmdFlags{FileA: "/nonexistent/a.env", FileB: "/nonexistent/b.env"}
	var buf bytes.Buffer
	err := Run(flags, &buf)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}
