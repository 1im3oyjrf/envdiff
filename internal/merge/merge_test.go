package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/merge"
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

func TestMerge_NoConflicts(t *testing.T) {
	src := writeTempEnv(t, "APP=myapp\nDEBUG=true\n")
	tgt := writeTempEnv(t, "PORT=8080\nHOST=localhost\n")

	res, err := merge.Merge(src, tgt, merge.StrategySource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Entries["APP"] != "myapp" || res.Entries["PORT"] != "8080" {
		t.Errorf("unexpected entries: %v", res.Entries)
	}
}

func TestMerge_ConflictStrategySource(t *testing.T) {
	src := writeTempEnv(t, "APP=source-app\n")
	tgt := writeTempEnv(t, "APP=target-app\n")

	res, err := merge.Merge(src, tgt, merge.StrategySource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Resolved != "source-app" {
		t.Errorf("expected source value, got %q", res.Conflicts[0].Resolved)
	}
	if res.Entries["APP"] != "source-app" {
		t.Errorf("expected source-app in entries, got %q", res.Entries["APP"])
	}
}

func TestMerge_ConflictStrategyTarget(t *testing.T) {
	src := writeTempEnv(t, "APP=source-app\n")
	tgt := writeTempEnv(t, "APP=target-app\n")

	res, err := merge.Merge(src, tgt, merge.StrategyTarget)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries["APP"] != "target-app" {
		t.Errorf("expected target-app, got %q", res.Entries["APP"])
	}
}

func TestMerge_InvalidSourceFile(t *testing.T) {
	tgt := writeTempEnv(t, "APP=x\n")
	_, err := merge.Merge("/nonexistent.env", tgt, merge.StrategySource)
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func TestWrite_OutputFile(t *testing.T) {
	src := writeTempEnv(t, "Z=last\nA=first\n")
	tgt := writeTempEnv(t, "M=middle\n")

	res, err := merge.Merge(src, tgt, merge.StrategySource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := filepath.Join(t.TempDir(), "merged.env")
	if err := merge.Write(res, out); err != nil {
		t.Fatalf("write error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	content := string(data)
	if content == "" {
		t.Error("expected non-empty output file")
	}
}
