package audit

import (
	"fmt"
	"sort"
	"strings"
)

// Severity represents the importance level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityError   Severity = "ERROR"
)

// Finding represents a single audit observation for a key.
type Finding struct {
	Key      string
	Severity Severity
	Message  string
}

// Result holds the full audit output for a set of env files.
type Result struct {
	Findings []Finding
	Total    int
	Errors   int
	Warnings int
	Infos    int
}

// AuditFiles runs an audit across the provided env maps and returns findings.
// source is the reference env, targets are additional envs keyed by label.
func AuditFiles(source map[string]string, targets map[string]map[string]string) Result {
	findings := []Finding{}

	for key, val := range source {
		if strings.TrimSpace(val) == "" {
			findings = append(findings, Finding{
				Key:      key,
				Severity: SeverityWarning,
				Message:  "empty value in source",
			})
		}
		for label, tenv := range targets {
			tval, ok := tenv[key]
			if !ok {
				findings = append(findings, Finding{
					Key:      key,
					Severity: SeverityError,
					Message:  fmt.Sprintf("missing in target %q", label),
				})
			} else if tval == val && strings.TrimSpace(val) != "" {
				findings = append(findings, Finding{
					Key:      key,
					Severity: SeverityInfo,
					Message:  fmt.Sprintf("identical value in target %q — consider environment-specific override", label),
				})
			}
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Key != findings[j].Key {
			return findings[i].Key < findings[j].Key
		}
		return findings[i].Message < findings[j].Message
	})

	result := Result{Findings: findings, Total: len(findings)}
	for _, f := range findings {
		switch f.Severity {
		case SeverityError:
			result.Errors++
		case SeverityWarning:
			result.Warnings++
		case SeverityInfo:
			result.Infos++
		}
	}
	return result
}
