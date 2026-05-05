package config

import (
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envdiff-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestValidate_MissingSource(t *testing.T) {
	cfg := &Config{SourceFile: "", TargetFile: "some.env", OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing source, got nil")
	}
}

func TestValidate_MissingTarget(t *testing.T) {
	src := writeTempFile(t, "KEY=value\n")
	cfg := &Config{SourceFile: src, TargetFile: "", OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing target, got nil")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	src := writeTempFile(t, "KEY=value\n")
	dst := writeTempFile(t, "KEY=value\n")
	cfg := &Config{SourceFile: src, TargetFile: dst, OutputFormat: "xml"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestValidate_SourceNotFound(t *testing.T) {
	dst := writeTempFile(t, "KEY=value\n")
	cfg := &Config{SourceFile: "/nonexistent/path.env", TargetFile: dst, OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}

func TestValidate_TargetNotFound(t *testing.T) {
	src := writeTempFile(t, "KEY=value\n")
	cfg := &Config{SourceFile: src, TargetFile: "/nonexistent/path.env", OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing target file, got nil")
	}
}

func TestValidate_Valid(t *testing.T) {
	src := writeTempFile(t, "KEY=value\n")
	dst := writeTempFile(t, "KEY=value\n")
	cfg := &Config{SourceFile: src, TargetFile: dst, OutputFormat: "json", MaskSecrets: true}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
