package template

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFlags_Valid(t *testing.T) {
	src := writeTempEnv(t, "APP=1\n")
	opts, err := ParseFlags([]string{"-output", "out.env", src}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Output != "out.env" {
		t.Errorf("expected output 'out.env', got %q", opts.Output)
	}
	if len(opts.Sources) != 1 || opts.Sources[0] != src {
		t.Errorf("unexpected sources: %v", opts.Sources)
	}
}

func TestParseFlags_DefaultOutput(t *testing.T) {
	src := writeTempEnv(t, "APP=1\n")
	opts, err := ParseFlags([]string{src}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Output != ".env.template" {
		t.Errorf("expected default output, got %q", opts.Output)
	}
}

func TestParseFlags_NoSources(t *testing.T) {
	_, err := ParseFlags([]string{}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error when no sources provided")
	}
}

func TestRun_WritesOutputAndPrints(t *testing.T) {
	src := writeTempEnv(t, "# app name\nAPP_NAME=myapp\nDB_HOST=localhost\n")
	out := filepath.Join(t.TempDir(), "result.env")

	opts := &Options{Sources: []string{src}, Output: out}
	var buf bytes.Buffer
	if err := Run(opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Generated template with 2 keys") {
		t.Errorf("unexpected stdout: %q", output)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "APP_NAME=") {
		t.Errorf("expected APP_NAME= in output, got:\n%s", data)
	}
}

func TestRun_MissingSource(t *testing.T) {
	opts := &Options{Sources: []string{"/no/such/file.env"}, Output: "/tmp/out.env"}
	err := Run(opts, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
