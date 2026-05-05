package diff

import (
	"testing"

	"github.com/yourusername/envdiff/internal/parser"
)

func TestCompare_Clean(t *testing.T) {
	src := parser.EnvMap{"KEY": "value", "PORT": "8080"}
	tgt := parser.EnvMap{"KEY": "value", "PORT": "8080"}

	result := Compare(src, tgt)
	if !result.IsClean() {
		t.Errorf("expected clean result, got %+v", result)
	}
}

func TestCompare_MissingInTarget(t *testing.T) {
	src := parser.EnvMap{"KEY": "value", "MISSING": "x"}
	tgt := parser.EnvMap{"KEY": "value"}

	result := Compare(src, tgt)
	if len(result.MissingInTarget) != 1 || result.MissingInTarget[0] != "MISSING" {
		t.Errorf("expected MISSING in MissingInTarget, got %v", result.MissingInTarget)
	}
}

func TestCompare_MissingInSource(t *testing.T) {
	src := parser.EnvMap{"KEY": "value"}
	tgt := parser.EnvMap{"KEY": "value", "EXTRA": "y"}

	result := Compare(src, tgt)
	if len(result.MissingInSource) != 1 || result.MissingInSource[0] != "EXTRA" {
		t.Errorf("expected EXTRA in MissingInSource, got %v", result.MissingInSource)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	src := parser.EnvMap{"KEY": "old"}
	tgt := parser.EnvMap{"KEY": "new"}

	result := Compare(src, tgt)
	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "KEY" || m.SourceValue != "old" || m.TargetValue != "new" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	src := parser.EnvMap{"Z_KEY": "v", "A_KEY": "v", "M_KEY": "v"}
	tgt := parser.EnvMap{}

	result := Compare(src, tgt)
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, key := range expected {
		if result.MissingInTarget[i] != key {
			t.Errorf("index %d: expected %q, got %q", i, key, result.MissingInTarget[i])
		}
	}
}

func TestIsClean_False(t *testing.T) {
	r := Result{MissingInTarget: []string{"KEY"}}
	if r.IsClean() {
		t.Error("expected IsClean to return false")
	}
}
