package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestWriteText_Clean(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{}
	err := WriteText(&buf, result, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestWriteText_MissingInTarget(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInTarget: []string{"DB_HOST", "DB_PORT"},
	}
	err := WriteText(&buf, result, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Missing in target (2)") {
		t.Errorf("expected missing in target header, got: %s", out)
	}
	if !strings.Contains(out, "- DB_HOST") {
		t.Errorf("expected DB_HOST entry, got: %s", out)
	}
}

func TestWriteText_Mismatched_MaskSecrets(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		Mismatched: []diff.Mismatch{
			{Key: "API_SECRET", SourceValue: "abc123", TargetValue: "xyz789"},
		},
	}
	opts := DefaultOptions()
	opts.MaskSecrets = true
	err := WriteText(&buf, result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "abc123") || strings.Contains(out, "xyz789") {
		t.Errorf("expected masked values, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected mask placeholder, got: %s", out)
	}
}

func TestWriteText_Summary(t *testing.T) {
	var buf bytes.Buffer
	result := diff.Result{
		MissingInTarget: []string{"FOO"},
		MissingInSource: []string{"BAR"},
	}
	err := WriteText(&buf, result, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "---") {
		t.Errorf("expected separator line, got: %s", out)
	}
}
