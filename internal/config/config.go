package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the parsed CLI configuration for envdiff.
type Config struct {
	SourceFile string
	TargetFile string
	MaskSecrets bool
	OutputFormat string
	ExitOnDiff bool
}

// Parse parses command-line flags and returns a Config.
func Parse() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.SourceFile, "source", "", "Path to the source .env file (required)")
	flag.StringVar(&cfg.TargetFile, "target", "", "Path to the target .env file (required)")
	flag.BoolVar(&cfg.MaskSecrets, "mask", false, "Mask secret values in output")
	flag.StringVar(&cfg.OutputFormat, "format", "text", "Output format: text or json")
	flag.BoolVar(&cfg.ExitOnDiff, "exit-on-diff", false, "Exit with code 1 if differences are found")

	flag.Parse()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required fields are set and values are valid.
func (c *Config) Validate() error {
	if c.SourceFile == "" {
		return fmt.Errorf("--source is required")
	}
	if c.TargetFile == "" {
		return fmt.Errorf("--target is required")
	}
	if c.OutputFormat != "text" && c.OutputFormat != "json" {
		return fmt.Errorf("--format must be \"text\" or \"json\", got %q", c.OutputFormat)
	}
	if _, err := os.Stat(c.SourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source file not found: %s", c.SourceFile)
	}
	if _, err := os.Stat(c.TargetFile); os.IsNotExist(err) {
		return fmt.Errorf("target file not found: %s", c.TargetFile)
	}
	return nil
}
