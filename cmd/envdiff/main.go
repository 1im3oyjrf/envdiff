package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/report"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Parse()
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	source, err := parser.ParseFile(cfg.SourceFile)
	if err != nil {
		return fmt.Errorf("parsing source file %q: %w", cfg.SourceFile, err)
	}

	target, err := parser.ParseFile(cfg.TargetFile)
	if err != nil {
		return fmt.Errorf("parsing target file %q: %w", cfg.TargetFile, err)
	}

	result := diff.Compare(source, target)

	opts := report.Options{
		MaskSecrets: cfg.MaskSecrets,
	}

	switch cfg.OutputFormat {
	case "json":
		if err := report.WriteJSON(os.Stdout, result, opts); err != nil {
			return fmt.Errorf("writing JSON report: %w", err)
		}
	default:
		if err := report.Write(os.Stdout, result, opts); err != nil {
			return fmt.Errorf("writing text report: %w", err)
		}
	}

	if cfg.ExitOnDiff && !result.Clean() {
		os.Exit(1)
	}

	return nil
}
