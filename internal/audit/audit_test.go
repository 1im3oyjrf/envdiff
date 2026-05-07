package audit

import (
	"testing"
)

func TestAuditFiles_Clean(t *testing.T) {
	source := map[string]string{"HOST": "localhost", "PORT": "8080"}
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com", "PORT": "443"},
	}
	result := AuditFiles(source, targets)
	if result.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", result.Errors)
	}
	if result.Warnings != 0 {
		t.Errorf("expected 0 warnings, got %d", result.Warnings)
	}
}

func TestAuditFiles_MissingInTarget(t *testing.T) {
	source := map[string]string{"HOST": "localhost", "SECRET": "abc"}
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com"},
	}
	result := AuditFiles(source, targets)
	if result.Errors != 1 {
		t.Errorf("expected 1 error, got %d", result.Errors)
	}
	found := false
	for _, f := range result.Findings {
		if f.Key == "SECRET" && f.Severity == SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected finding for missing key SECRET")
	}
}

func TestAuditFiles_EmptyValueInSource(t *testing.T) {
	source := map[string]string{"HOST": "", "PORT": "8080"}
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.example.com", "PORT": "443"},
	}
	result := AuditFiles(source, targets)
	if result.Warnings != 1 {
		t.Errorf("expected 1 warning, got %d", result.Warnings)
	}
}

func TestAuditFiles_IdenticalValues(t *testing.T) {
	source := map[string]string{"HOST": "same", "PORT": "9000"}
	targets := map[string]map[string]string{
		"staging": {"HOST": "same", "PORT": "9000"},
	}
	result := AuditFiles(source, targets)
	if result.Infos != 2 {
		t.Errorf("expected 2 infos for identical values, got %d", result.Infos)
	}
}

func TestAuditFiles_SortedFindings(t *testing.T) {
	source := map[string]string{"ZEBRA": "", "APPLE": ""}
	targets := map[string]map[string]string{}
	result := AuditFiles(source, targets)
	if len(result.Findings) < 2 {
		t.Fatal("expected at least 2 findings")
	}
	if result.Findings[0].Key > result.Findings[1].Key {
		t.Error("findings are not sorted by key")
	}
}

func TestAuditFiles_MultipleTargets(t *testing.T) {
	source := map[string]string{"KEY": "val"}
	targets := map[string]map[string]string{
		"staging": {},
		"prod":    {},
	}
	result := AuditFiles(source, targets)
	if result.Errors != 2 {
		t.Errorf("expected 2 errors (one per target), got %d", result.Errors)
	}
}
