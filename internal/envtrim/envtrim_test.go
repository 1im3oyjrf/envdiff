package envtrim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestTrimFile_NoChangesNeeded(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	res, err := TrimFile(p, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != 0 {
		t.Errorf("expected 0 trimmed, got %d", res.Trimmed)
	}
	if res.Total != 2 {
		t.Errorf("expected 2 total, got %d", res.Total)
	}
}

func TestTrimFile_TrimsWhitespace(t *testing.T) {
	p := writeTempEnv(t, "KEY=  hello  \nFOO=\tbar\t\nCLEAN=ok\n")
	res, err := TrimFile(p, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != 2 {
		t.Errorf("expected 2 trimmed, got %d", res.Trimmed)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "KEY=hello") {
		t.Errorf("expected KEY=hello in output, got: %s", string(data))
	}
	if !strings.Contains(string(data), "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", string(data))
	}
}

func TestTrimFile_DryRunDoesNotWrite(t *testing.T) {
	original := "KEY=  spaced  \n"
	p := writeTempEnv(t, original)
	res, err := TrimFile(p, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun to be true")
	}
	if res.Trimmed != 1 {
		t.Errorf("expected 1 trimmed, got %d", res.Trimmed)
	}
	// File should be unchanged.
	data, _ := os.ReadFile(p)
	if string(data) != original {
		t.Errorf("expected file unchanged in dry-run, got: %s", string(data))
	}
}

func TestTrimFile_PreservesCommentsAndBlanks(t *testing.T) {
	input := "# comment\n\nKEY=value\n"
	p := writeTempEnv(t, input)
	_, err := TrimFile(p, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "# comment") {
		t.Errorf("comment should be preserved")
	}
}

func TestTrimFile_WritesToOutputFile(t *testing.T) {
	p := writeTempEnv(t, "KEY=  val  \n")
	outDir := t.TempDir()
	out := filepath.Join(outDir, "out.env")
	_, err := TrimFile(p, out, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "KEY=val") {
		t.Errorf("expected KEY=val in output file, got: %s", string(data))
	}
	// Source should be unchanged.
	src, _ := os.ReadFile(p)
	if !strings.Contains(string(src), "KEY=  val  ") {
		t.Errorf("source should be unchanged")
	}
}

func TestTrimFile_MissingFile(t *testing.T) {
	_, err := TrimFile("/nonexistent/.env", "", false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
