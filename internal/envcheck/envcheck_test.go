package envcheck

import (
	"testing"
)

func TestCheckAgainstEnv_AllMatch(t *testing.T) {
	t.Setenv("APP_HOST", "localhost")
	t.Setenv("APP_PORT", "8080")

	envMap := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}

	results := CheckAgainstEnv(envMap, Options{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Present || !r.Match {
			t.Errorf("expected key %s to be present and matching", r.Key)
		}
	}
}

func TestCheckAgainstEnv_MissingKey(t *testing.T) {
	envMap := map[string]string{
		"DEFINITELY_NOT_SET_XYZ": "value",
	}

	results := CheckAgainstEnv(envMap, Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Present {
		t.Error("expected key to be absent from environment")
	}
	if results[0].Match {
		t.Error("absent key should not match")
	}
}

func TestCheckAgainstEnv_Mismatch(t *testing.T) {
	t.Setenv("APP_ENV", "production")

	envMap := map[string]string{
		"APP_ENV": "development",
	}

	results := CheckAgainstEnv(envMap, Options{})
	if results[0].Match {
		t.Error("expected mismatch")
	}
	if !results[0].Present {
		t.Error("expected key to be present")
	}
}

func TestCheckAgainstEnv_IgnoreCase(t *testing.T) {
	t.Setenv("APP_MODE", "DEBUG")

	envMap := map[string]string{
		"APP_MODE": "debug",
	}

	results := CheckAgainstEnv(envMap, Options{IgnoreCase: true})
	if !results[0].Match {
		t.Error("expected case-insensitive match")
	}
}

func TestSummary_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Present: true, Match: true},
		{Key: "B", Present: false, Match: false},
		{Key: "C", Present: true, Match: false},
	}

	present, missing, mismatched := Summary(results)
	if present != 1 || missing != 1 || mismatched != 1 {
		t.Errorf("unexpected counts: present=%d missing=%d mismatched=%d", present, missing, mismatched)
	}
}

func TestFormatMismatch_Missing(t *testing.T) {
	r := Result{Key: "SECRET_KEY", Expected: "abc123", Present: false}
	out := FormatMismatch(r, false)
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestFormatMismatch_MaskSecrets(t *testing.T) {
	r := Result{Key: "DB_PASS", Expected: "hunter2", Actual: "wrong", Present: true, Match: false}
	out := FormatMismatch(r, true)
	if contains(out, "hunter2") || contains(out, "wrong") {
		t.Error("secret values should be masked")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
