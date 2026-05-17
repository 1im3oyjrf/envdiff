package envhide_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envhide"
)

func TestParseFlags_Valid(t *testing.T) {
	f, err := envhide.ParseFlags([]string{"-source", ".env", "-placeholder", "HIDDEN"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Source != ".env" {
		t.Errorf("expected source .env, got %q", f.Source)
	}
	if f.Placeholder != "HIDDEN" {
		t.Errorf("expected placeholder HIDDEN, got %q", f.Placeholder)
	}
}

func TestParseFlags_MissingSource(t *testing.T) {
	_, err := envhide.ParseFlags([]string{})
	if err == nil {
		t.Error("expected error for missing -source")
	}
	if !strings.Contains(err.Error(), "-source") {
		t.Errorf("error should mention -source, got: %v", err)
	}
}

func TestParseFlags_DefaultPlaceholder(t *testing.T) {
	f, err := envhide.ParseFlags([]string{"-source", ".env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Placeholder != "***" {
		t.Errorf("expected default placeholder ***, got %q", f.Placeholder)
	}
}

func TestRun_WritesToFile(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=myapp\nSECRET_KEY=topsecret\n")
	out := filepath.Join(t.TempDir(), "out.env")

	f := envhide.Flags{Source: src, Output: out, Placeholder: "REDACTED"}
	if err := envhide.Run(f); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if strings.Contains(string(data), "topsecret") {
		t.Error("output should not contain original secret value")
	}
	if !strings.Contains(string(data), "REDACTED") {
		t.Error("output should contain placeholder")
	}
}

func TestRun_MissingSource(t *testing.T) {
	f := envhide.Flags{Source: "/no/such/.env", Placeholder: "***"}
	if err := envhide.Run(f); err == nil {
		t.Error("expected error for missing source file")
	}
}
