package envfilter_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envfilter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envfilter-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestFilterFile_Prefix(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\nAPP_ENV=production\n")
	r, err := envfilter.FilterFile(path, envfilter.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if r.Matched["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", r.Matched["DB_HOST"])
	}
	if r.Total != 4 {
		t.Errorf("expected total 4, got %d", r.Total)
	}
}

func TestFilterFile_Suffix(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nAPP_HOST=example.com\nAPP_PORT=8080\n")
	r, err := envfilter.FilterFile(path, envfilter.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilterFile_KeyPattern(t *testing.T) {
	path := writeTempEnv(t, "SECRET_KEY=abc\nAPI_SECRET=xyz\nAPP_NAME=test\n")
	r, err := envfilter.FilterFile(path, envfilter.Options{KeyPattern: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilterFile_Invert(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")
	r, err := envfilter.FilterFile(path, envfilter.Options{Prefix: "DB_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 1 {
		t.Errorf("expected 1 matched after invert, got %d", len(r.Matched))
	}
	if _, ok := r.Matched["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in inverted results")
	}
}

func TestFilterFile_MissingFile(t *testing.T) {
	_, err := envfilter.FilterFile("/nonexistent/path.env", envfilter.Options{})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWriteText_Output(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")
	r, err := envfilter.FilterFile(path, envfilter.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	envfilter.WriteText(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2/3 matched") {
		t.Errorf("expected match count in output, got:\n%s", out)
	}
}
