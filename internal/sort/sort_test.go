package sort

import (
	"bytes"
	"os"
	"strings"
	"testing"
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

func TestSortFile_Alphabetical(t *testing.T) {
	src := writeTempEnv(t, "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")
	var buf bytes.Buffer
	err := SortFile(src, &buf, Options{Alphabetical: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmpty(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APPLE") {
		t.Errorf("expected APPLE first, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "MIDDLE") {
		t.Errorf("expected MIDDLE second, got %q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected ZEBRA last, got %q", lines[2])
	}
}

func TestSortFile_PreservesAllKeys(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\nHELLO=world\n")
	var buf bytes.Buffer
	err := SortFile(src, &buf, Options{Alphabetical: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"FOO", "BAZ", "HELLO"} {
		if !strings.Contains(out, key) {
			t.Errorf("output missing key %q", key)
		}
	}
}

func TestSortFile_MissingSource(t *testing.T) {
	var buf bytes.Buffer
	err := SortFile("/no/such/file.env", &buf, Options{Alphabetical: true})
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestSortFile_EmptyFile(t *testing.T) {
	src := writeTempEnv(t, "")
	var buf bytes.Buffer
	err := SortFile(src, &buf, Options{Alphabetical: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

// nonEmpty returns non-blank lines from s.
func nonEmpty(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
