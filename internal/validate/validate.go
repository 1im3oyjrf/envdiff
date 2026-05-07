package validate

import (
	"fmt"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Allowed  []string // if non-empty, value must be one of these
	NoEmpty  bool
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Message string
}

// Result holds all violations found during validation.
type Result struct {
	Violations []Violation
}

// OK returns true if no violations were found.
func (r *Result) OK() bool {
	return len(r.Violations) == 0
}

// CheckEnv validates a parsed env map against a set of rules.
func CheckEnv(env map[string]string, rules []Rule) *Result {
	result := &Result{}

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.NoEmpty && strings.TrimSpace(val) == "" {
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: "value must not be empty",
			})
		}

		if len(rule.Allowed) > 0 && !contains(rule.Allowed, val) {
			result.Violations = append(result.Violations, Violation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value %q is not allowed; must be one of: %s", val, strings.Join(rule.Allowed, ", ")),
			})
		}
	}

	return result
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
