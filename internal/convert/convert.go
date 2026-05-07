package convert

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON    Format = "json"
	FormatYAML    Format = "yaml"
)

// Result holds the converted output string and the detected format.
type Result struct {
	Format Format
	Output string
}

// Convert transforms a key-value map into the specified output format.
func Convert(env map[string]string, format Format) (Result, error) {
	keys := sortedKeys(env)

	switch format {
	case FormatDotenv:
		return Result{Format: format, Output: toDotenv(env, keys)}, nil
	case FormatExport:
		return Result{Format: format, Output: toExport(env, keys)}, nil
	case FormatJSON:
		return Result{Format: format, Output: toJSON(env, keys)}, nil
	case FormatYAML:
		return Result{Format: format, Output: toYAML(env, keys)}, nil
	default:
		return Result{}, fmt.Errorf("unsupported format: %q", format)
	}
}

func toDotenv(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func toExport(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func toJSON(env map[string]string, keys []string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, env[k], comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func toYAML(env map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, ":\n#") {
			fmt.Fprintf(&sb, "%s: %q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s: %s\n", k, v)
		}
	}
	return sb.String()
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
