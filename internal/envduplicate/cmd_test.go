package envduplicate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvCmd(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestParseFlags_Valid(t *testing.T) {
	f, err := ParseFlags([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(f.Files))
	}
}

func TestParseFlags_SummaryFlag(t *testing.T) {
	f, err := ParseFlags([]string{"-summary", "a.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Summary {
		t.Error("expected Summary=true")
	}
}

func TestParseFlags_NoFiles(t *testing.T) {
	_, err := ParseFlags([]string{})
	if err == nil {
		t.Error("expected error when no files given")
	}
}

func TestRun_WritesToStdout(t *testing.T) {
	p := writeTempEnvCmd(t, "A=dup\nB=dup\n")
	f := &Flags{Files: []string{p}}
	var sb strings.Builder
	if err := Run(f, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "dup") {
		t.Errorf("expected 'dup' in output: %s", sb.String())
	}
}

func TestRun_SummaryMode(t *testing.T) {
	p := writeTempEnvCmd(t, "X=v\nY=v\n")
	f := &Flags{Files: []string{p}, Summary: true}
	var sb strings.Builder
	if err := Run(f, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "duplicate-value group") {
		t.Errorf("expected summary line in output: %s", sb.String())
	}
}

func TestRun_WritesToFile(t *testing.T) {
	p := writeTempEnvCmd(t, "A=same\nB=same\n")
	outFile := filepath.Join(t.TempDir(), "out.txt")
	f := &Flags{Files: []string{p}, Output: outFile}
	var sb strings.Builder
	if err := Run(f, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "same") {
		t.Errorf("expected 'same' in file output")
	}
}
