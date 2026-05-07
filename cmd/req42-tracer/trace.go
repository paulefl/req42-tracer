package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/internal/model"
	"github.com/paulefl/req42-tracer/internal/parser"
	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/report"
)

func newTraceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Generate traceability matrix",
		Long: `Generate and display the traceability matrix showing relationships
between requirements, architecture, and tests.`,
		RunE: runTraceCmd,
	}

	return cmd
}

func runTraceCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	format, _ := cmd.Flags().GetString("format")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load configuration (for future use with advanced options)
	_, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	// Build traceability graph
	builder := graph.NewBuilder()

	// Parse requirements from docs/requirements/
	if req, err := parser.ParseAllFromDir("docs/requirements", "software"); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements from docs/requirements/\n")
		}
	}

	// Parse architecture from docs/arc42/
	if arch, err := parser.ParseAllFromDir("docs/arc42", "software"); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture from docs/arc42/\n")
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

	// Generate report
	analyzer := graph.NewAnalyzer(g)
	reporter := report.NewTableReporter(analyzer, format)
	output := reporter.TraceabilityMatrix()

	fmt.Print(output)

	return nil
}
