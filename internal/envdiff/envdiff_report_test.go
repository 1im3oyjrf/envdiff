package envdiff

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteText_NoDifferences(t *testing.T) {
	result := &DiffResult{
		FileA: "a.env",
		FileB: "b.env",
	}
	var buf bytes.Buffer
	WriteText(&buf, result, false)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteText_OnlyInA(t *testing.T) {
	result := &DiffResult{
		FileA:  "a.env",
		FileB:  "b.env",
		OnlyInA: map[string]string{"FOO": "bar"},
	}
	var buf bytes.Buffer
	WriteText(&buf, result, false)
	out := buf.String()
	if !strings.Contains(out, "Only in a.env") {
		t.Errorf("expected section header, got: %s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar, got: %s", out)
	}
}

func TestWriteText_MaskSecrets(t *testing.T) {
	result := &DiffResult{
		FileA:  "a.env",
		FileB:  "b.env",
		OnlyInA: map[string]string{"SECRET_KEY": "supersecret"},
	}
	var buf bytes.Buffer
	WriteText(&buf, result, true)
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("secret value should be masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected masked value, got: %s", out)
	}
}

func TestWriteText_Mismatched(t *testing.T) {
	result := &DiffResult{
		FileA: "a.env",
		FileB: "b.env",
		Mismatched: []MismatchedEntry{
			{Key: "DB_HOST", ValueA: "localhost", ValueB: "prod.db"},
		},
	}
	var buf bytes.Buffer
	WriteText(&buf, result, false)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "prod.db") {
		t.Errorf("expected both values in output, got: %s", out)
	}
}

func TestWriteSummary_Counts(t *testing.T) {
	result := &DiffResult{
		FileA:   "a.env",
		FileB:   "b.env",
		OnlyInA: map[string]string{"X": "1"},
		OnlyInB: map[string]string{"Y": "2", "Z": "3"},
		Mismatched: []MismatchedEntry{
			{Key: "W", ValueA: "a", ValueB: "b"},
		},
	}
	var buf bytes.Buffer
	WriteSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "1 only-in-A") {
		t.Errorf("expected 1 only-in-A, got: %s", out)
	}
	if !strings.Contains(out, "2 only-in-B") {
		t.Errorf("expected 2 only-in-B, got: %s", out)
	}
	if !strings.Contains(out, "1 mismatched") {
		t.Errorf("expected 1 mismatched, got: %s", out)
	}
}
