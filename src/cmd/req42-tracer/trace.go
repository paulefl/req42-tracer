package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
	"github.com/paulefl/req42-tracer/src/internal/report"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

func newTraceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Generate traceability matrix",
		Long: `Generate and display the traceability matrix showing relationships
between requirements, architecture, and tests. Supports text, markdown, JSON,
and interactive HTML visualization with D3.js graph.`,
		RunE: runTraceCmd,
	}

	cmd.Flags().String("output", "", "Output file path for HTML report (optional)")

	return cmd
}

func runTraceCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	format, _ := cmd.Flags().GetString("format")
	outputPath, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	// Build traceability graph
	builder := graph.NewBuilder()

	// Parse requirements from project/req42-tracer/docs/requirements/
	if req, err := parser.ParseAllFromDir("project/req42-tracer/docs/requirements", "software"); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements from project/req42-tracer/docs/requirements/\n")
		}
	}

	// Parse architecture from project/req42-tracer/docs/arc42/
	if arch, err := parser.ParseAllFromDir("project/req42-tracer/docs/arc42", "software"); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture from project/req42-tracer/docs/arc42/\n")
		}
	}

	// Load Bausteinsicht model if configured
	if bPath := config.Bausteinsicht.Model; bPath != "" {
		bParser := parser.NewBausteinsichtParser(bPath)
		if bGraph, err := bParser.Parse("software"); err == nil {
			if err := builder.MergeGraph(bGraph); err != nil {
				return err
			}
			if verbose {
				fmt.Fprintf(os.Stderr, "Parsed Bausteinsicht model from %s\n", bPath)
			}
		} else if verbose {
			fmt.Fprintf(os.Stderr, "Warning: could not load Bausteinsicht model: %v\n", err)
		}
	}

	// Derive ASPICE levels
	builder.DeriveASPICELevels()

	// Build trace links
	if err := builder.BuildLinks(); err != nil {
		return err
	}

	// Get final graph
	g := builder.GetGraph()

	// Load test results from CI artifacts (JUnit XML, go-test JSON)
	if err := testresult.LoadAll(g, config); err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: could not load test results: %v\n", err)
	}

	// Generate report
	analyzer := graph.NewAnalyzer(g)

	// Handle HTML output separately
	if format == "html" || outputPath != "" {
		if outputPath == "" {
			outputPath = "reports/graph.html"
		}
		htmlReporter := report.NewHTMLReporter(analyzer, config, outputPath)
		if err := htmlReporter.GenerateReport(); err != nil {
			return fmt.Errorf("failed to generate HTML report: %w", err)
		}
		summaryPath := filepath.Join(filepath.Dir(outputPath), "summary.html")
		if err := htmlReporter.GenerateSummaryReport(summaryPath); err != nil {
			return fmt.Errorf("failed to generate summary report: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "HTML report generated: %s\n", outputPath)
			fmt.Fprintf(os.Stderr, "Summary report generated: %s\n", summaryPath)
		}
		fmt.Printf("HTML reports generated:\n  Graph: %s\n  Summary: %s\n", outputPath, summaryPath)
		return nil
	}

	// Use table reporter for text/markdown/json formats
	reporter := report.NewTableReporter(analyzer, format)
	output := reporter.TraceabilityMatrix()

	fmt.Print(output)

	return nil
}
