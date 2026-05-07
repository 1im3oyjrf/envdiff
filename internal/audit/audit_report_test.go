package audit

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteText_Clean(t *testing.T) {
	result := Result{}
	var buf bytes.Buffer
	WriteText(&buf, result)
	if !strings.Contains(buf.String(), "Audit passed") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestWriteText_WithFindings(t *testing.T) {
	result := Result{
		Findings: []Finding{
			{Key: "DB_PASS", Severity: SeverityError, Message: `missing in target "prod"`},
			{Key: "HOST", Severity: SeverityWarning, Message: "empty value in source"},
		},
		Total:    2,
		Errors:   1,
		Warnings: 1,
	}
	var buf bytes.Buffer
	WriteText(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "DB_PASS") {
		t.Error("expected DB_PASS in output")
	}
	if !strings.Contains(out, "HOST") {
		t.Error("expected HOST in output")
	}
	if !strings.Contains(out, "Summary") {
		t.Error("expected Summary line")
	}
}

func TestWriteSummary_Counts(t *testing.T) {
	result := Result{Total: 5, Errors: 2, Warnings: 2, Infos: 1}
	var buf bytes.Buffer
	WriteSummary(&buf, result)
	out := buf.String()
	if !strings.Contains(out, "5 finding(s)") {
		t.Errorf("expected total count in summary, got: %s", out)
	}
	if !strings.Contains(out, "2 error(s)") {
		t.Errorf("expected error count in summary, got: %s", out)
	}
}

func TestWriteText_IconsPresent(t *testing.T) {
	result := Result{
		Findings: []Finding{
			{Key: "A", Severity: SeverityError, Message: "err"},
			{Key: "B", Severity: SeverityWarning, Message: "warn"},
			{Key: "C", Severity: SeverityInfo, Message: "info"},
		},
		Total: 3, Errors: 1, Warnings: 1, Infos: 1,
	}
	var buf bytes.Buffer
	WriteText(&buf, result)
	out := buf.String()
	for _, icon := range []string{"✖", "⚠", "ℹ"} {
		if !strings.Contains(out, icon) {
			t.Errorf("expected icon %q in output", icon)
		}
	}
}
