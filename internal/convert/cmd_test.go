package convert

import (
	"bytes"
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

func TestParseFlags_Valid(t *testing.T) {
	f, err := ParseFlags([]string{"--source", "a.env", "--format", "json", "--output", "out.json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Source != "a.env" || f.Format != "json" || f.Output != "out.json" {
		t.Errorf("unexpected flags: %+v", f)
	}
}

func TestParseFlags_MissingSource(t *testing.T) {
	_, err := ParseFlags([]string{"--format", "yaml"})
	if err == nil {
		t.Error("expected error for missing --source")
	}
}

func TestParseFlags_DefaultFormat(t *testing.T) {
	f, err := ParseFlags([]string{"--source", "x.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Format != "dotenv" {
		t.Errorf("expected default format dotenv, got %s", f.Format)
	}
}

func TestRun_StdoutOutput(t *testing.T) {
	src := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	var buf bytes.Buffer
	err := Run(Flags{Source: src, Format: "dotenv"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY=value") {
		t.Errorf("expected output to contain KEY=value, got: %s", buf.String())
	}
}

func TestRun_FileOutput(t *testing.T) {
	src := writeTempEnv(t, "APP=test\n")
	out := filepath.Join(t.TempDir(), "out.json")
	var buf bytes.Buffer
	err := Run(Flags{Source: src, Format: "json", Output: out}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), `"APP"`) {
		t.Errorf("expected JSON output file to contain APP key")
	}
	if !strings.Contains(buf.String(), "Converted") {
		t.Errorf("expected confirmation message in stdout")
	}
}

func TestRun_MissingSource(t *testing.T) {
	err := Run(Flags{Source: "nonexistent.env", Format: "dotenv"}, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
