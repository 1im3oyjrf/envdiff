package copy

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFlags_Valid(t *testing.T) {
	flags, err := ParseFlags([]string{"--source", "a.env", "--dest", "b.env", "--key", "FOO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flags.Source != "a.env" || flags.Dest != "b.env" || flags.Key != "FOO" {
		t.Errorf("unexpected flags: %+v", flags)
	}
}

func TestParseFlags_MissingSource(t *testing.T) {
	_, err := ParseFlags([]string{"--dest", "b.env", "--key", "FOO"})
	if err == nil || !strings.Contains(err.Error(), "--source") {
		t.Errorf("expected --source error, got: %v", err)
	}
}

func TestParseFlags_MissingDest(t *testing.T) {
	_, err := ParseFlags([]string{"--source", "a.env", "--key", "FOO"})
	if err == nil || !strings.Contains(err.Error(), "--dest") {
		t.Errorf("expected --dest error, got: %v", err)
	}
}

func TestParseFlags_MissingKey(t *testing.T) {
	_, err := ParseFlags([]string{"--source", "a.env", "--dest", "b.env"})
	if err == nil || !strings.Contains(err.Error(), "--key") {
		t.Errorf("expected --key error, got: %v", err)
	}
}

func TestRun_AddsKey(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dest := filepath.Join(dir, "dest.env")
	os.WriteFile(src, []byte("TOKEN=secret\n"), 0644)
	os.WriteFile(dest, []byte("OTHER=val\n"), 0644)

	var buf bytes.Buffer
	err := Run([]string{"--source", src, "--dest", dest, "--key", "TOKEN"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "added") {
		t.Errorf("expected 'added' in output, got: %s", buf.String())
	}
	data, _ := os.ReadFile(dest)
	if !strings.Contains(string(data), "TOKEN=secret") {
		t.Errorf("expected TOKEN in dest, got: %s", data)
	}
}

func TestRun_UpdatesKey(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.env")
	dest := filepath.Join(dir, "dest.env")
	os.WriteFile(src, []byte("TOKEN=newval\n"), 0644)
	os.WriteFile(dest, []byte("TOKEN=oldval\n"), 0644)

	var buf bytes.Buffer
	err := Run([]string{"--source", src, "--dest", dest, "--key", "TOKEN", "--overwrite"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "updated") {
		t.Errorf("expected 'updated' in output, got: %s", buf.String())
	}
}

func TestRun_MissingSourceFile(t *testing.T) {
	var buf bytes.Buffer
	err := Run([]string{"--source", "/no/such/file.env", "--dest", "/tmp/d.env", "--key", "X"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}
