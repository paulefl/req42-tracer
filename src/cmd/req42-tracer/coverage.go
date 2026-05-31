package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/report"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

func newCoverageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coverage",
		Short: "Generate coverage dashboard from test coverage data",
		Long: `Generate an interactive coverage dashboard that combines statement-level
test coverage data with the traceability graph.

Supported input formats:
  --coverage    go tool cover -func output (Go)
  --cobertura   Cobertura XML (Java, Python, .NET, C/C++)
  --lcov        LCOV .info file (C/C++, JavaScript)

The dashboard maps packages to architecture elements via the impl= attribute
and shows coverage per component and per package.

Examples:
  # Go project
  go test ./... -coverprofile=coverage.out
  req42-tracer coverage --coverage coverage.out --output coverage.html

  # Cobertura (e.g. Maven)
  req42-tracer coverage --cobertura target/site/cobertura/coverage.xml

  # LCOV (e.g. gcov)
  req42-tracer coverage --lcov lcov.info --output coverage.html`,
		RunE: runCoverageCmd,
	}

	cmd.Flags().String("coverage", "", "go tool cover -func output file")
	cmd.Flags().String("cobertura", "", "Cobertura XML coverage file")
	cmd.Flags().String("lcov", "", "LCOV .info coverage file")
	cmd.Flags().String("output", "coverage.html", "Output HTML file path")

	return cmd
}

func runCoverageCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	gocover, _ := cmd.Flags().GetString("coverage")
	cobertura, _ := cmd.Flags().GetString("cobertura")
	lcov, _ := cmd.Flags().GetString("lcov")
	output, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	if gocover == "" && cobertura == "" && lcov == "" {
		return fmt.Errorf("at least one coverage source required: --coverage, --cobertura, or --lcov")
	}

	// Load config + graph for arch mapping
	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	g, err := buildGraph(config,
		"project/req42-tracer/docs/requirements",
		"project/req42-tracer/docs/arc42",
		config.GetDefaultProject(), verbose)
	if err != nil {
		return fmt.Errorf("failed to build graph: %w", err)
	}

	// Parse coverage sources
	var allPkgs []testresult.PackageCoverage

	if gocover != "" {
		pkgs, err := testresult.ParseGoCoverage(gocover)
		if err != nil {
			return fmt.Errorf("go coverage: %w", err)
		}
		allPkgs = append(allPkgs, pkgs...)
		if verbose {
			fmt.Fprintf(os.Stderr, "Loaded %d packages from go coverage\n", len(pkgs))
		}
	}

	if cobertura != "" {
		pkgs, err := testresult.ParseCobertura(cobertura)
		if err != nil {
			return fmt.Errorf("cobertura: %w", err)
		}
		allPkgs = append(allPkgs, pkgs...)
		if verbose {
			fmt.Fprintf(os.Stderr, "Loaded %d packages from Cobertura\n", len(pkgs))
		}
	}

	if lcov != "" {
		pkgs, err := testresult.ParseLCOV(lcov)
		if err != nil {
			return fmt.Errorf("lcov: %w", err)
		}
		allPkgs = append(allPkgs, pkgs...)
		if verbose {
			fmt.Fprintf(os.Stderr, "Loaded %d packages from LCOV\n", len(pkgs))
		}
	}

	covData := report.BuildCoverageData(allPkgs, g)

	if strings.HasSuffix(output, ".json") {
		data, _ := json.MarshalIndent(covData, "", "  ")
		return os.WriteFile(output, data, 0644)
	}

	// Generate HTML report
	covJSON, err := json.MarshalIndent(covData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal coverage data: %w", err)
	}

	htmlContent := strings.ReplaceAll(report.CoverageHTMLTemplate, "<!--COVERAGE_DATA_JSON-->", string(covJSON))

	dir := filepath.Dir(output)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	if err := os.WriteFile(output, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("write coverage report: %w", err)
	}

	fmt.Printf("Coverage report written to %s\n", output)
	fmt.Printf("Overall: %.1f%% (%d/%d statements)\n",
		covData.OverallPct, covData.TotalCov, covData.TotalStmts)
	return nil
}
