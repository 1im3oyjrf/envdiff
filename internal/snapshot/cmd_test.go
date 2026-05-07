package snapshot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvSnap(t *testing.T, content string) string {
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

func TestParseFlags_CaptureAction(t *testing.T) {
	f, err := ParseFlags([]string{"--env", ".env", "capture"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Action != "capture" {
		t.Errorf("expected action capture, got %s", f.Action)
	}
	if f.SnapshotFile != ".env.snapshot" {
		t.Errorf("expected default snapshot file, got %s", f.SnapshotFile)
	}
}

func TestParseFlags_MissingEnv(t *testing.T) {
	_, err := ParseFlags([]string{"capture"})
	if err == nil || !strings.Contains(err.Error(), "--env") {
		t.Errorf("expected --env error, got %v", err)
	}
}

func TestParseFlags_UnknownAction(t *testing.T) {
	_, err := ParseFlags([]string{"--env", ".env", "delete"})
	if err == nil || !strings.Contains(err.Error(), "unknown action") {
		t.Errorf("expected unknown action error, got %v", err)
	}
}

func TestParseFlags_MissingAction(t *testing.T) {
	_, err := ParseFlags([]string{"--env", ".env"})
	if err == nil || !strings.Contains(err.Error(), "action required") {
		t.Errorf("expected action required error, got %v", err)
	}
}

func TestRun_CaptureAndCompare(t *testing.T) {
	envContent := "APP_ENV=production\nDB_HOST=localhost\n"
	envFile := writeTempEnvSnap(t, envContent)
	snapshotFile := filepath.Join(t.TempDir(), "snap.json")

	// Capture
	flags := &Flags{
		Action:       "capture",
		EnvFile:      envFile,
		SnapshotFile: snapshotFile,
		MaskSecrets:  false,
	}
	if err := Run(flags, os.Stdout); err != nil {
		t.Fatalf("capture failed: %v", err)
	}

	// Compare (no changes)
	flags.Action = "compare"
	if err := Run(flags, os.Stdout); err != nil {
		t.Fatalf("compare failed: %v", err)
	}
}

func TestRun_CompareMissingSnapshot(t *testing.T) {
	envFile := writeTempEnvSnap(t, "KEY=val\n")
	flags := &Flags{
		Action:       "compare",
		EnvFile:      envFile,
		SnapshotFile: "/nonexistent/snap.json",
	}
	if err := Run(flags, os.Stdout); err == nil {
		t.Error("expected error for missing snapshot file")
	}
}
