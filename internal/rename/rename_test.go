package rename

import (
	"os"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRenameKey_Success(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	res, err := RenameKey(path, "DB_HOST", "DATABASE_HOST", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Errorf("expected Renamed=true, got false")
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "DATABASE_HOST=localhost") {
		t.Errorf("renamed key not found in file: %s", string(raw))
	}
	if strings.Contains(string(raw), "DB_HOST=") {
		t.Errorf("old key still present in file")
	}
}

func TestRenameKey_NotFound(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\n")
	res, err := RenameKey(path, "MISSING_KEY", "NEW_KEY", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Errorf("expected Skipped=true")
	}
	if res.Reason != "key not found" {
		t.Errorf("unexpected reason: %s", res.Reason)
	}
}

func TestRenameKey_TargetExists_NoOverwrite(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=foo\nNEW_KEY=bar\n")
	res, err := RenameKey(path, "OLD_KEY", "NEW_KEY", Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Skipped {
		t.Errorf("expected Skipped=true when target exists and overwrite=false")
	}
}

func TestRenameKey_TargetExists_WithOverwrite(t *testing.T) {
	path := writeTempEnv(t, "OLD_KEY=foo\nNEW_KEY=bar\n")
	res, err := RenameKey(path, "OLD_KEY", "NEW_KEY", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Errorf("expected Renamed=true with overwrite=true")
	}
}

func TestRenameKey_DryRun(t *testing.T) {
	original := "DB_NAME=mydb\n"
	path := writeTempEnv(t, original)
	res, err := RenameKey(path, "DB_NAME", "DATABASE_NAME", Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Renamed {
		t.Errorf("expected Renamed=true in dry run")
	}
	raw, _ := os.ReadFile(path)
	if string(raw) != original {
		t.Errorf("file should be unchanged in dry run")
	}
}

func TestRenameKey_EmptyKeys(t *testing.T) {
	path := writeTempEnv(t, "KEY=val\n")
	_, err := RenameKey(path, "", "NEW", Options{})
	if err == nil {
		t.Error("expected error for empty oldKey")
	}
}
