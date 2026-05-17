package envdiff

import (
	"os"
	"path/filepath"
	"testing"
)

func makeFileSummary(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDiff_NoChanges(t *testing.T) {
	a := makeFileSummary(t, "FOO=bar\nBAZ=qux\n")
	b := makeFileSummary(t, "FOO=bar\nBAZ=qux\n")
	res, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInA)+len(res.OnlyInB)+len(res.Mismatched) != 0 {
		t.Errorf("expected no diff, got %+v", res)
	}
}

func TestDiff_OnlyInA(t *testing.T) {
	a := makeFileSummary(t, "FOO=bar\nEXTRA=yes\n")
	b := makeFileSummary(t, "FOO=bar\n")
	res, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.OnlyInA["EXTRA"]; !ok {
		t.Error("expected EXTRA in OnlyInA")
	}
}

func TestDiff_OnlyInB(t *testing.T) {
	a := makeFileSummary(t, "FOO=bar\n")
	b := makeFileSummary(t, "FOO=bar\nNEW=val\n")
	res, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.OnlyInB["NEW"]; !ok {
		t.Error("expected NEW in OnlyInB")
	}
}

func TestDiff_Mismatched(t *testing.T) {
	a := makeFileSummary(t, "DB=localhost\n")
	b := makeFileSummary(t, "DB=prod\n")
	res, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Mismatched) != 1 || res.Mismatched[0].Key != "DB" {
		t.Errorf("expected DB mismatch, got %+v", res.Mismatched)
	}
}

func TestDiff_SecretKeyDetected(t *testing.T) {
	if !isSecretKey("API_SECRET") {
		t.Error("expected API_SECRET to be a secret key")
	}
	if isSecretKey("APP_NAME") {
		t.Error("expected APP_NAME not to be a secret key")
	}
}

func TestDiff_MissingFileReturnsError(t *testing.T) {
	_, err := Diff("/nonexistent/a.env", "/nonexistent/b.env")
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}

func TestDiff_FileNamesPreserved(t *testing.T) {
	a := makeFileSummary(t, "X=1\n")
	b := makeFileSummary(t, "X=1\n")
	res, err := Diff(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if res.FileA != a || res.FileB != b {
		t.Errorf("file names not preserved: %s %s", res.FileA, res.FileB)
	}
}
