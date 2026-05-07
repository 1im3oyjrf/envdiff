package validate

import (
	"testing"
)

func TestCheckEnv_AllValid(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"LOG_LEVEL": "info",
	}
	rules := []Rule{
		{Key: "APP_ENV", Required: true, Allowed: []string{"production", "staging", "development"}},
		{Key: "LOG_LEVEL", Required: true, NoEmpty: true},
	}
	result := CheckEnv(env, rules)
	if !result.OK() {
		t.Errorf("expected no violations, got: %v", result.Violations)
	}
}

func TestCheckEnv_MissingRequired(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{
		{Key: "DATABASE_URL", Required: true},
	}
	result := CheckEnv(env, rules)
	if result.OK() {
		t.Fatal("expected violation for missing required key")
	}
	if result.Violations[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key: %s", result.Violations[0].Key)
	}
}

func TestCheckEnv_EmptyValueNotAllowed(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "   "}
	rules := []Rule{
		{Key: "SECRET_KEY", NoEmpty: true},
	}
	result := CheckEnv(env, rules)
	if result.OK() {
		t.Fatal("expected violation for empty value")
	}
	if result.Violations[0].Message != "value must not be empty" {
		t.Errorf("unexpected message: %s", result.Violations[0].Message)
	}
}

func TestCheckEnv_DisallowedValue(t *testing.T) {
	env := map[string]string{"APP_ENV": "test"}
	rules := []Rule{
		{Key: "APP_ENV", Allowed: []string{"production", "staging"}},
	}
	result := CheckEnv(env, rules)
	if result.OK() {
		t.Fatal("expected violation for disallowed value")
	}
}

func TestCheckEnv_OptionalMissingKeySkipped(t *testing.T) {
	env := map[string]string{}
	rules := []Rule{
		{Key: "OPTIONAL_FLAG", Allowed: []string{"true", "false"}},
	}
	result := CheckEnv(env, rules)
	if !result.OK() {
		t.Errorf("expected no violations for missing optional key, got: %v", result.Violations)
	}
}

func TestCheckEnv_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "unknown",
	}
	rules := []Rule{
		{Key: "APP_ENV", Allowed: []string{"production"}},
		{Key: "DB_HOST", Required: true},
	}
	result := CheckEnv(env, rules)
	if len(result.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(result.Violations))
	}
}
