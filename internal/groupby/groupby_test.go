package groupby

import (
	"os"
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

func TestGroupByPrefix_Basic(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAWS_KEY=abc\nPORT=8080\n")
	groups, err := GroupByPrefix(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if groups[0].Name != "AWS" {
		t.Errorf("expected first group AWS, got %s", groups[0].Name)
	}
	if groups[1].Name != "DB" {
		t.Errorf("expected second group DB, got %s", groups[1].Name)
	}
	if groups[2].Name != "OTHER" {
		t.Errorf("expected third group OTHER, got %s", groups[2].Name)
	}
}

func TestGroupByPrefix_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nDB_HOST=localhost\n")
	groups, err := GroupByPrefix(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != "DB" {
		t.Errorf("expected 1 DB group, got %+v", groups)
	}
}

func TestGroupByPrefix_SortedEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_USER=root\nDB_HOST=localhost\nDB_PORT=5432\n")
	groups, err := GroupByPrefix(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group")
	}
	keys := []string{groups[0].Entries[0].Key, groups[0].Entries[1].Key, groups[0].Entries[2].Key}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PORT" || keys[2] != "DB_USER" {
		t.Errorf("entries not sorted: %v", keys)
	}
}

func TestGroupByPrefix_MissingFile(t *testing.T) {
	_, err := GroupByPrefix("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestGroupByPrefix_NoUnderscore(t *testing.T) {
	path := writeTempEnv(t, "PORT=3000\nHOST=localhost\n")
	groups, err := GroupByPrefix(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != "OTHER" {
		t.Errorf("expected OTHER group, got %+v", groups)
	}
	if len(groups[0].Entries) != 2 {
		t.Errorf("expected 2 entries in OTHER, got %d", len(groups[0].Entries))
	}
}
