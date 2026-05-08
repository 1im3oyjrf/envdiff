package envcheck

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteText_AllOK(t *testing.T) {
	results := []Result{
		{Key: "APP_HOST", Present: true, Match: true},
		{Key: "APP_PORT", Present: true, Match: true},
	}

	var buf bytes.Buffer
	WriteText(&buf, results, false)
	out := buf.String()

	if !strings.Contains(out, "OK") {
		t.Error("expected OK entries in output")
	}
	if !strings.Contains(out, "2 ok") {
		t.Error("expected summary to show 2 ok")
	}
}

func TestWriteText_WithMissing(t *testing.T) {
	results := []Result{
		{Key: "DB_URL", Expected: "postgres://localhost", Present: false},
	}

	var buf bytes.Buffer
	WriteText(&buf, results, false)
	out := buf.String()

	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING label")
	}
	if !strings.Contains(out, "1 missing") {
		t.Error("expected summary to show 1 missing")
	}
}

func TestWriteText_MaskSecrets(t *testing.T) {
	results := []Result{
		{Key: "DB_PASSWORD", Expected: "s3cr3t", Actual: "wrong", Present: true, Match: false},
	}

	var buf bytes.Buffer
	WriteText(&buf, results, true)
	out := buf.String()

	if strings.Contains(out, "s3cr3t") || strings.Contains(out, "wrong") {
		t.Error("secret values should not appear in masked output")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked placeholder in output")
	}
}

func TestWriteSummaryText(t *testing.T) {
	results := []Result{
		{Key: "A", Present: true, Match: true},
		{Key: "B", Present: false},
		{Key: "C", Present: true, Match: false},
	}

	var buf bytes.Buffer
	WriteSummaryText(&buf, results)
	out := buf.String()

	if !strings.Contains(out, "1 ok") {
		t.Errorf("expected '1 ok' in %q", out)
	}
	if !strings.Contains(out, "1 missing") {
		t.Errorf("expected '1 missing' in %q", out)
	}
	if !strings.Contains(out, "1 mismatched") {
		t.Errorf("expected '1 mismatched' in %q", out)
	}
}

func TestHasIssues(t *testing.T) {
	clean := []Result{{Key: "X", Present: true, Match: true}}
	if HasIssues(clean) {
		t.Error("expected no issues for clean results")
	}

	dirty := []Result{{Key: "Y", Present: false}}
	if !HasIssues(dirty) {
		t.Error("expected issues for missing key")
	}
}
