package envdiff

import (
	"bytes"
	"strings"
	"testing"
)

func makeFileSummary(path string, entries map[string]string) FileSummary {
	return FileSummary{Path: path, Entries: entries}
}

func TestDiff_NoChanges(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"FOO": "bar", "BAZ": "qux"})
	b := makeFileSummary("b.env", map[string]string{"FOO": "bar", "BAZ": "qux"})

	result := Diff(a, b)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got: %+v", result)
	}
	if len(result.Identical) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(result.Identical))
	}
}

func TestDiff_OnlyInA(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"FOO": "bar", "EXTRA": "value"})
	b := makeFileSummary("b.env", map[string]string{"FOO": "bar"})

	result := Diff(a, b)

	if len(result.OnlyInA) != 1 || result.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA in OnlyInA, got %v", result.OnlyInA)
	}
}

func TestDiff_OnlyInB(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"FOO": "bar"})
	b := makeFileSummary("b.env", map[string]string{"FOO": "bar", "NEW": "val"})

	result := Diff(a, b)

	if len(result.OnlyInB) != 1 || result.OnlyInB[0] != "NEW" {
		t.Errorf("expected NEW in OnlyInB, got %v", result.OnlyInB)
	}
}

func TestDiff_Mismatched(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"FOO": "bar"})
	b := makeFileSummary("b.env", map[string]string{"FOO": "baz"})

	result := Diff(a, b)

	if len(result.Different) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Different))
	}
	if result.Different[0].Key != "FOO" || result.Different[0].ValueA != "bar" || result.Different[0].ValueB != "baz" {
		t.Errorf("unexpected mismatch: %+v", result.Different[0])
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"Z": "1", "A": "2", "M": "3"})
	b := makeFileSummary("b.env", map[string]string{})

	result := Diff(a, b)

	for i := 1; i < len(result.OnlyInA); i++ {
		if result.OnlyInA[i] < result.OnlyInA[i-1] {
			t.Errorf("OnlyInA not sorted: %v", result.OnlyInA)
		}
	}
}

func TestWriteText_NoDiff(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"FOO": "bar"})
	b := makeFileSummary("b.env", map[string]string{"FOO": "bar"})
	result := Diff(a, b)

	var buf bytes.Buffer
	WriteText(&buf, a, b, result, false)

	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected clean output, got: %s", buf.String())
	}
}

func TestWriteText_MaskSecrets(t *testing.T) {
	a := makeFileSummary("a.env", map[string]string{"API_SECRET": "abc123"})
	b := makeFileSummary("b.env", map[string]string{"API_SECRET": "xyz789"})
	result := Diff(a, b)

	var buf bytes.Buffer
	WriteText(&buf, a, b, result, true)

	output := buf.String()
	if strings.Contains(output, "abc123") || strings.Contains(output, "xyz789") {
		t.Errorf("expected secrets to be masked, got: %s", output)
	}
	if !strings.Contains(output, "***") {
		t.Errorf("expected masked placeholder in output, got: %s", output)
	}
}
