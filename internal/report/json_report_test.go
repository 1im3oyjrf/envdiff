package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func TestWriteJSON_Clean(t *testing.T) {
	result := diff.Result{}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, result, DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalIssues != 0 {
		t.Errorf("expected 0 issues, got %d", out.TotalIssues)
	}
	if len(out.MissingKeys) != 0 {
		t.Errorf("expected empty missing_keys")
	}
	if len(out.ExtraKeys) != 0 {
		t.Errorf("expected empty extra_keys")
	}
}

func TestWriteJSON_WithDiffs(t *testing.T) {
	result := diff.Result{
		MissingInTarget: []string{"DB_HOST"},
		MissingInSource: []string{"OLD_KEY"},
		Mismatched: []diff.Mismatch{
			{Key: "APP_ENV", SourceValue: "production", TargetValue: "staging"},
		},
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, result, DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalIssues != 3 {
		t.Errorf("expected 3 issues, got %d", out.TotalIssues)
	}
	if len(out.MissingKeys) != 1 || out.MissingKeys[0] != "DB_HOST" {
		t.Errorf("unexpected missing_keys: %v", out.MissingKeys)
	}
	if len(out.Mismatched) != 1 || out.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("unexpected mismatched: %v", out.Mismatched)
	}
}

func TestWriteJSON_MaskSecrets(t *testing.T) {
	result := diff.Result{
		Mismatched: []diff.Mismatch{
			{Key: "SECRET_KEY", SourceValue: "abc123", TargetValue: "xyz789"},
		},
	}
	opts := Options{MaskSecrets: true}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, result, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out JSONReport
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Mismatched[0].SourceValue != "***" {
		t.Errorf("expected masked source value, got %q", out.Mismatched[0].SourceValue)
	}
	if out.Mismatched[0].TargetValue != "***" {
		t.Errorf("expected masked target value, got %q", out.Mismatched[0].TargetValue)
	}
}
