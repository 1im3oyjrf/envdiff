package envhide_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envhide"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestHideSecrets_SecretsAreHidden(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nSECRET_KEY=abc123\nDB_PASSWORD=hunter2\n")
	res, err := envhide.HideSecrets(p, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HiddenCount != 2 {
		t.Errorf("expected 2 hidden, got %d", res.HiddenCount)
	}
	joined := strings.Join(res.Lines, "\n")
	if strings.Contains(joined, "abc123") {
		t.Error("secret value abc123 should be hidden")
	}
	if strings.Contains(joined, "hunter2") {
		t.Error("secret value hunter2 should be hidden")
	}
	if !strings.Contains(joined, "APP_NAME=myapp") {
		t.Error("non-secret APP_NAME should be preserved")
	}
}

func TestHideSecrets_CommentsAndBlanksPreserved(t *testing.T) {
	p := writeTempEnv(t, "# comment\n\nAPP_ENV=production\n")
	res, err := envhide.HideSecrets(p, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HiddenCount != 0 {
		t.Errorf("expected 0 hidden, got %d", res.HiddenCount)
	}
	if res.Lines[0] != "# comment" {
		t.Errorf("expected comment preserved, got %q", res.Lines[0])
	}
	if res.Lines[1] != "" {
		t.Errorf("expected blank line preserved, got %q", res.Lines[1])
	}
}

func TestHideSecrets_DefaultPlaceholder(t *testing.T) {
	p := writeTempEnv(t, "API_TOKEN=supersecret\n")
	res, err := envhide.HideSecrets(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(strings.Join(res.Lines, "\n"), "***") {
		t.Error("expected default placeholder *** to be used")
	}
}

func TestHideSecrets_MissingFile(t *testing.T) {
	_, err := envhide.HideSecrets("/nonexistent/.env", "***")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWrite_ToFile(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nSECRET_KEY=s3cr3t\n")
	res, err := envhide.HideSecrets(p, "REDACTED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := filepath.Join(t.TempDir(), "out.env")
	if err := envhide.Write(res, out); err != nil {
		t.Fatalf("write error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "s3cr3t") {
		t.Error("output file should not contain original secret")
	}
	if !strings.Contains(string(data), "REDACTED") {
		t.Error("output file should contain placeholder")
	}
}
