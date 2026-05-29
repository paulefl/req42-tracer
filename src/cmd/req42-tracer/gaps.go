package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/report"
)

func newGapsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gaps",
		Short: "Analyze gaps and orphan artifacts",
		Long: `Analyze traceability gaps including orphan requirements,
architecture elements without requirements, and missing implementations.`,
		RunE: runGapsCmd,
	}

	return cmd
}

func runGapsCmd(cmd *cobra.Command, args []string) error {
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

	// Parse requirements
	if req, err := parser.ParseAllFromDir("project/req42-tracer/docs/requirements", "software"); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements\n")
		}
	}

	// Parse architecture
	if arch, err := parser.ParseAllFromDir("project/req42-tracer/docs/arc42", "software"); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture\n")
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

	// Analyze gaps
	analyzer := graph.NewAnalyzer(g)

	// Generate report
	reporter := report.NewTableReporter(analyzer, format)
	output := reporter.GapReport()

	fmt.Print(output)

	return nil
}
