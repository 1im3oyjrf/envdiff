package redact_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/redact"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRedactFile_SecretsAreRedacted(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=myapp\nSECRET_KEY=supersecret\nDB_PASSWORD=hunter2\n")
	var buf bytes.Buffer
	res, err := redact.RedactFile(src, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Error("expected SECRET_KEY value to be redacted")
	}
	if strings.Contains(out, "hunter2") {
		t.Error("expected DB_PASSWORD value to be redacted")
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Error("expected non-secret APP_NAME to be preserved")
	}
	if res.LinesRedacted != 2 {
		t.Errorf("expected 2 redacted lines, got %d", res.LinesRedacted)
	}
}

func TestRedactFile_CommentsAndBlanksPreserved(t *testing.T) {
	src := writeTempEnv(t, "# comment\n\nAPP_ENV=production\n")
	var buf bytes.Buffer
	res, err := redact.RedactFile(src, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "# comment") {
		t.Error("expected comment to be preserved")
	}
	if res.LinesRedacted != 0 {
		t.Errorf("expected 0 redacted lines, got %d", res.LinesRedacted)
	}
}

func TestRedactFile_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	_, err := redact.RedactFile("/nonexistent/.env", &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
