package envduplicate

import (
	"strings"
	"testing"
)

func TestWriteText_NoDuplicates(t *testing.T) {
	results := []Result{
		{File: "clean.env", Findings: nil},
	}
	var sb strings.Builder
	WriteText(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "No duplicate values found") {
		t.Errorf("expected clean message, got: %s", out)
	}
}

func TestWriteText_MultipleDuplicates(t *testing.T) {
	results := []Result{
		{
			File: "multi.env",
			Findings: []Finding{
				{Value: "abc", Keys: []string{"KEY1", "KEY2"}},
				{Value: "xyz", Keys: []string{"FOO", "BAR", "BAZ"}},
			},
		},
	}
	var sb strings.Builder
	WriteText(&sb, results)
	out := sb.String()
	for _, want := range []string{"abc", "KEY1", "KEY2", "xyz", "FOO", "BAR", "BAZ"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func TestWriteSummary_ZeroFindings(t *testing.T) {
	results := []Result{
		{File: "empty.env", Findings: []Finding{}},
	}
	var sb strings.Builder
	WriteSummary(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "0 duplicate-value group") {
		t.Errorf("expected zero count, got: %s", out)
	}
}

func TestWriteSummary_MultipleResults(t *testing.T) {
	results := []Result{
		{File: "a.env", Findings: []Finding{
			{Value: "v1", Keys: []string{"K1", "K2"}},
			{Value: "v2", Keys: []string{"K3", "K4"}},
		}},
		{File: "b.env", Findings: nil},
	}
	var sb strings.Builder
	WriteSummary(&sb, results)
	out := sb.String()
	if !strings.Contains(out, "2 duplicate-value group") {
		t.Errorf("expected count 2 in summary: %s", out)
	}
	if !strings.Contains(out, "b.env: 0") {
		t.Errorf("expected b.env with 0 in summary: %s", out)
	}
}

func TestDetectDuplicates_SortedKeys(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "shared",
		"ALPHA": "shared",
		"MANGO": "shared",
	}
	findings := detectDuplicates(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	keys := findings[0].Keys
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}
