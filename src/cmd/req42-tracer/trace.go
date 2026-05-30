package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
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

	g, err := buildGraph(config,
		"project/req42-tracer/docs/requirements",
		"project/req42-tracer/docs/arc42",
		config.GetDefaultProject(), verbose)
	if err != nil {
		return err
	}

	if err := testresult.LoadAll(g, config); err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: could not load test results: %v\n", err)
	}

	analyzer := graph.NewAnalyzer(g)

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

	reporter := report.NewTableReporter(analyzer, format)
	fmt.Print(reporter.TraceabilityMatrix())
	return nil
}
