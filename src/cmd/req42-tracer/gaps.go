package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
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

	analyzer := graph.NewAnalyzer(g)
	fmt.Print(report.NewTableReporter(analyzer, format).GapReport())
	return nil
}
