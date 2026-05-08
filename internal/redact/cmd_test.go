package redact_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/redact"
)

func TestParseFlags_Valid(t *testing.T) {
	f, err := redact.ParseFlags([]string{"path/to/.env", "--output", "out.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Source != "path/to/.env" {
		t.Errorf("expected source path/to/.env, got %s", f.Source)
	}
	if f.Output != "out.env" {
		t.Errorf("expected output out.env, got %s", f.Output)
	}
}

func TestParseFlags_NoSource(t *testing.T) {
	_, err := redact.ParseFlags([]string{})
	if err == nil {
		t.Fatal("expected error when source is missing")
	}
}

func TestRun_WritesToStdout(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=test\nAPI_SECRET=abc123\n")
	var buf bytes.Buffer
	if err := redact.Run([]string{src}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "abc123") {
		t.Error("expected secret to be redacted in output")
	}
	if !strings.Contains(out, "Redacted") {
		t.Error("expected summary line in output")
	}
}

func TestRun_WritesToFile(t *testing.T) {
	src := writeTempEnv(t, "DB_PASSWORD=secret\nAPP_ENV=prod\n")
	outPath := filepath.Join(t.TempDir(), "redacted.env")
	var buf bytes.Buffer
	if err := redact.Run([]string{src, "--output", outPath}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if strings.Contains(string(data), "secret") {
		t.Error("expected DB_PASSWORD to be redacted in output file")
	}
}
