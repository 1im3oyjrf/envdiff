package deletekey

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
		t.Fatalf("writeTempEnvCmd: %v", err)
	}
	return p
}

func TestParseFlags_Valid(t *testing.T) {
	f, err := ParseFlags([]string{"--file", ".env", "--key", "FOO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.File != ".env" || f.Key != "FOO" || f.DryRun {
		t.Errorf("unexpected flags: %+v", f)
	}
}

func TestParseFlags_MissingFile(t *testing.T) {
	_, err := ParseFlags([]string{"--key", "FOO"})
	if err == nil || !strings.Contains(err.Error(), "--file") {
		t.Errorf("expected --file error, got: %v", err)
	}
}

func TestParseFlags_MissingKey(t *testing.T) {
	_, err := ParseFlags([]string{"--file", ".env"})
	if err == nil || !strings.Contains(err.Error(), "--key") {
		t.Errorf("expected --key error, got: %v", err)
	}
}

func TestParseFlags_DryRun(t *testing.T) {
	f, err := ParseFlags([]string{"--file", ".env", "--key", "FOO", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.DryRun {
		t.Error("expected DryRun=true")
	}
}

func TestRun_DeletesKey(t *testing.T) {
	p := writeTempEnvCmd(t, "FOO=bar\nSECRET=xyz\n")
	var sb strings.Builder
	err := Run(Flags{File: p, Key: "SECRET"}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "deleted") {
		t.Errorf("expected 'deleted' in output, got: %q", sb.String())
	}
}

func TestRun_KeyNotFound(t *testing.T) {
	p := writeTempEnvCmd(t, "FOO=bar\n")
	var sb strings.Builder
	err := Run(Flags{File: p, Key: "MISSING"}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "not found") {
		t.Errorf("expected 'not found' in output, got: %q", sb.String())
	}
}

func TestRun_DryRunOutput(t *testing.T) {
	p := writeTempEnvCmd(t, "FOO=bar\nSECRET=xyz\n")
	var sb strings.Builder
	err := Run(Flags{File: p, Key: "SECRET", DryRun: true}, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "dry-run") {
		t.Errorf("expected 'dry-run' in output, got: %q", sb.String())
	}
}
