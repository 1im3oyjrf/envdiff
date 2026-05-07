package convert

import (
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_NAME": "myapp",
	"DB_HOST":  "localhost",
	"SECRET":   "abc123",
}

func TestConvert_Dotenv(t *testing.T) {
	res, err := Convert(sampleEnv, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Format != FormatDotenv {
		t.Errorf("expected format dotenv, got %s", res.Format)
	}
	for _, line := range []string{"APP_NAME=myapp", "DB_HOST=localhost", "SECRET=abc123"} {
		if !strings.Contains(res.Output, line) {
			t.Errorf("expected output to contain %q", line)
		}
	}
}

func TestConvert_Export(t *testing.T) {
	res, err := Convert(sampleEnv, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "export APP_NAME=") {
		t.Errorf("expected export prefix in output")
	}
}

func TestConvert_JSON(t *testing.T) {
	res, err := Convert(sampleEnv, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(res.Output, "{") {
		t.Errorf("expected JSON output to start with {")
	}
	if !strings.Contains(res.Output, `"APP_NAME": "myapp"`) {
		t.Errorf("expected JSON to contain APP_NAME key")
	}
}

func TestConvert_YAML(t *testing.T) {
	res, err := Convert(sampleEnv, FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, "APP_NAME: myapp") {
		t.Errorf("expected YAML output to contain APP_NAME")
	}
}

func TestConvert_YAML_SpecialChars(t *testing.T) {
	env := map[string]string{"URL": "http://host:8080"}
	res, err := Convert(env, FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(res.Output, `URL: "http://host:8080"`) {
		t.Errorf("expected quoted YAML value for special chars, got: %s", res.Output)
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	_, err := Convert(sampleEnv, Format("toml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestConvert_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	res, err := Convert(env, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(res.Output), "\n")
	if lines[0] != "A_KEY=a" || lines[1] != "M_KEY=m" || lines[2] != "Z_KEY=z" {
		t.Errorf("output not sorted: %v", lines)
	}
}
