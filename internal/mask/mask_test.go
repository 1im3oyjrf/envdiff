package mask

import "testing"

func TestIsSecret_KnownSecretKeys(t *testing.T) {
	secretKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"APP_SECRET",
		"AWS_SECRET_ACCESS_KEY",
		"GITHUB_TOKEN",
		"USER_PASSWD",
	}
	for _, key := range secretKeys {
		if !IsSecret(key) {
			t.Errorf("expected IsSecret(%q) = true", key)
		}
	}
}

func TestIsSecret_NonSecretKeys(t *testing.T) {
	plainKeys := []string{
		"APP_ENV",
		"PORT",
		"DEBUG",
		"LOG_LEVEL",
		"DATABASE_URL",
	}
	for _, key := range plainKeys {
		if IsSecret(key) {
			t.Errorf("expected IsSecret(%q) = false", key)
		}
	}
}

func TestMaskValue(t *testing.T) {
	if got := MaskValue("supersecret"); got != "****" {
		t.Errorf("MaskValue: got %q, want %q", got, "****")
	}
	if got := MaskValue(""); got != "" {
		t.Errorf("MaskValue empty: got %q, want %q", got, "")
	}
}

func TestApplyMask_Enabled(t *testing.T) {
	got := ApplyMask("DB_PASSWORD", "s3cr3t", true)
	if got != "****" {
		t.Errorf("ApplyMask enabled secret: got %q, want ****", got)
	}
}

func TestApplyMask_Disabled(t *testing.T) {
	got := ApplyMask("DB_PASSWORD", "s3cr3t", false)
	if got != "s3cr3t" {
		t.Errorf("ApplyMask disabled: got %q, want s3cr3t", got)
	}
}

func TestApplyMask_NonSecretEnabled(t *testing.T) {
	got := ApplyMask("APP_ENV", "production", true)
	if got != "production" {
		t.Errorf("ApplyMask non-secret: got %q, want production", got)
	}
}
