package envduplicate

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

func TestFindDuplicateValues_NoDuplicates(t *testing.T) {
	f := writeTempEnv(t, "A=alpha\nB=beta\nC=gamma\n")
	results, err := FindDuplicateValues([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Findings) != 0 {
		t.Errorf("expected no findings, got %v", results[0].Findings)
	}
}

func TestFindDuplicateValues_WithDuplicates(t *testing.T) {
	f := writeTempEnv(t, "A=same\nB=same\nC=different\n")
	results, err := FindDuplicateValues([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(results[0].Findings))
	}
	f0 := results[0].Findings[0]
	if f0.Value != "same" {
		t.Errorf("expected value 'same', got %q", f0.Value)
	}
	if len(f0.Keys) != 2 {
		t.Errorf("expected 2 keys, got %v", f0.Keys)
	}
}

func TestFindDuplicateValues_IgnoresEmptyValues(t *testing.T) {
	f := writeTempEnv(t, "A=\nB=\nC=value\n")
	results, err := FindDuplicateValues([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Findings) != 0 {
		t.Errorf("empty values should not be flagged as duplicates")
	}
}

func TestFindDuplicateValues_MissingFile(t *testing.T) {
	_, err := FindDuplicateValues([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFindDuplicateValues_MultipleFiles(t *testing.T) {
	f1 := writeTempEnv(t, "X=dup\nY=dup\n")
	f2 := writeTempEnv(t, "P=unique\nQ=also_unique\n")
	results, err := FindDuplicateValues([]string{f1, f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Findings) != 1 {
		t.Errorf("file1: expected 1 finding")
	}
	if len(results[1].Findings) != 0 {
		t.Errorf("file2: expected no findings")
	}
}

func TestWriteText_Output(t *testing.T) {
	results := []Result{
		{
			File: "test.env",
			Findings: []Finding{
				{Value: "secret", Keys: []string{"API_KEY", "TOKEN"}},
			},
		},
	}
	var sb strings.Builder
	WriteText(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "secret") {
		t.Errorf("expected value in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected key in output")
	}
}

func TestWriteSummary_Output(t *testing.T) {
	results := []Result{
		{File: "a.env", Findings: []Finding{{Value: "v", Keys: []string{"K1", "K2"}}}},
		{File: "b.env", Findings: nil},
	}
	var sb strings.Builder
	WriteSummary(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "1 duplicate-value group") {
		t.Errorf("expected count in summary: %s", out)
	}
	if !strings.Contains(out, "0 duplicate-value group") {
		t.Errorf("expected zero count in summary: %s", out)
	}
}
