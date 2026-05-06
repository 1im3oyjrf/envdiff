package lint

import (
	"testing"
)

func TestCheckFile_Clean(t *testing.T) {
	lines := []string{
		"# This is a comment",
		"APP_ENV=production",
		"PORT=8080",
		"",
		"DB_HOST=localhost",
	}
	result := CheckFile("test.env", lines)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(result.Issues), result.Issues)
	}
	if result.HasErrors() {
		t.Error("expected HasErrors() to be false")
	}
}

func TestCheckFile_DuplicateKey(t *testing.T) {
	lines := []string{
		"APP_ENV=staging",
		"APP_ENV=production",
	}
	result := CheckFile("test.env", lines)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "warn" {
		t.Errorf("expected severity 'warn', got %q", result.Issues[0].Severity)
	}
	if result.Issues[0].Key != "APP_ENV" {
		t.Errorf("expected key 'APP_ENV', got %q", result.Issues[0].Key)
	}
}

func TestCheckFile_EmptyValue(t *testing.T) {
	lines := []string{
		"SECRET_KEY=",
		"API_TOKEN=\"\"",
	}
	result := CheckFile("test.env", lines)
	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %+v", len(result.Issues), result.Issues)
	}
	for _, issue := range result.Issues {
		if issue.Severity != "warn" {
			t.Errorf("expected severity 'warn', got %q", issue.Severity)
		}
	}
}

func TestCheckFile_InvalidKeyFormat(t *testing.T) {
	lines := []string{
		"INVALID-KEY=value",
		"123START=value",
	}
	result := CheckFile("test.env", lines)
	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(result.Issues))
	}
	for _, issue := range result.Issues {
		if issue.Severity != "error" {
			t.Errorf("expected severity 'error', got %q", issue.Severity)
		}
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors() to be true")
	}
}

func TestCheckFile_MissingEquals(t *testing.T) {
	lines := []string{
		"BADLINE",
	}
	result := CheckFile("test.env", lines)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %q", result.Issues[0].Severity)
	}
}

func TestIsValidKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"APP_ENV", true},
		{"port", true},
		{"DB_HOST_1", true},
		{"", false},
		{"INVALID-KEY", false},
		{"HAS SPACE", false},
		{"HAS.DOT", false},
	}
	for _, tc := range cases {
		got := isValidKey(tc.key)
		if got != tc.valid {
			t.Errorf("isValidKey(%q) = %v, want %v", tc.key, got, tc.valid)
		}
	}
}
